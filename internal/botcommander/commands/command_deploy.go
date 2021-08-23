package commands

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/unanet/eve-bot/internal/chatservice/chatmodels"
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/args"
	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
)

type deployCmd struct {
	baseCommand
}

const (
	// DeployCmdName is used as key/id for the deploy command
	DeployCmdName = "deploy"
)

var (
	deployCmdHelpSummary = help.Summary("The `deploy` command is used to deploy services to a specific *namespace* and *environment*")
	deployCmdHelpUsage   = help.Usage{
		"deploy {{ namespace }} in {{ environment }}",
		"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
		"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
		"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
	}
	deployCmdHelpExample = help.Examples{
		"deploy current in int",
		"deploy current in int services=api dryrun=true",
		"deploy current in int services=api,billing dryrun=true force=true",
		"deploy current in int services=api:1.0,billing",
	}
)

// NewDeployCommand creates a New DeployCmd that implements the EvebotCommand interface
func NewDeployCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := deployCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: DeployCmdName},
		arguments:  args.Args{args.DefaultDryrunArg(), args.DefaultForceArg(), args.DefaultServicesArg()},
		parameters: params.Params{params.DefaultNamespace(), params.DefaultEnvironment()},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 4, Max: 7},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd deployCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(deployCmdHelpSummary.String()),
		help.UsageOpt(deployCmdHelpUsage.String()),
		help.ArgsOpt(cmd.arguments.String()),
		help.ExamplesOpt(deployCmdHelpExample.String()),
	).String())
}

func (cmd deployCmd)  IsAuthenticated(chatUser *chatmodels.ChatUser, db *dynamodb.DynamoDB) bool {
	return true
}

// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd deployCmd) IsAuthorized(allowedChannel map[string]interface{}, chatChanFn ChatChannelInfoFn, chatUserFn ChatUserInfoFn, db *dynamodb.DynamoDB) bool {
	return cmd.IsHelpRequest() ||
		validChannelAuthCheck(cmd.info.Channel, allowedChannel, chatChanFn) ||
		lowerEnvAuthCheck(cmd.opts)
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd deployCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd deployCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *deployCmd) resolveDynamicOptions() {
	cmd.verifyInput()
	if len(cmd.errs) > 0 {
		return
	}

	cmd.opts[params.NamespaceName] = cmd.input[1]
	cmd.opts[params.EnvironmentName] = cmd.input[3]

	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if suppliedArg := args.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.opts[suppliedArg.Name()] = suppliedArg.Value()
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", argKV))
			}
		}
	}
}
