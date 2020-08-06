package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

func NewDeployCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultDeployCommand(cmdFields, channel, user)
}

type DeployCmd struct {
	baseCommand
}

func defaultDeployCommand(cmdFields []string, channel, user string) DeployCmd {
	cmd := DeployCmd{baseCommand{
		input:       cmdFields,
		chatDetails: ChatDetails{User: user, Channel: channel},
		name:        "deploy",
		summary:     "The `deploy` command is used to deploy services to a specific *namespace* and *environment*",
		usage: help.Usage{
			"deploy {{ namespace }} in {{ environment }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		examples: help.Examples{
			"deploy current in una-int",
			"deploy current in una-int services=unanetbi dryrun=true",
			"deploy current in una-int services=unanetbi,unaneta dryrun=true force=true",
			"deploy current in una-int services=unanetbi:20.2,unaneta",
		},
		optionalArgs:   args.Args{args.DefaultDryrunArg(), args.DefaultForceArg(), args.DefaultServicesArg()},
		requiredParams: params.Params{params.DefaultNamespace(), params.DefaultEnvironment()},
		apiOptions:     make(CommandOptions),
		inputBounds:    InputLengthBounds{Min: 4, Max: 7},
	}}
	cmd.resolveParams()
	cmd.resolveArgs()
	return cmd
}

func (cmd DeployCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfo) bool {
	return validChannelAuthCheck(cmd.chatDetails.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.apiOptions)
}

func (cmd DeployCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd DeployCmd) ChatInfo() ChatDetails {
	return cmd.chatDetails
}

func (cmd DeployCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd DeployCmd) IsValid() bool {
	return cmd.ValidInputLength()
}

func (cmd DeployCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd DeployCmd) Name() string {
	return cmd.name
}

func (cmd DeployCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ArgsOpt(cmd.optionalArgs.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd DeployCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

// resolveParams attempts to resolve the input params
// be sure and use a pointer receiver here since we are modifying the receiver object
func (cmd *DeployCmd) resolveParams() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("resolve cmd params err invalid input: %v", cmd.input))
		return
	}
	cmd.apiOptions[params.NamespaceName] = cmd.input[1]
	cmd.apiOptions[params.EnvironmentName] = cmd.input[3]
}

// resolveArgs attempts to resolve the input argument
// be sure and use a pointer receiver here since we are modifying the receiver object
func (cmd *DeployCmd) resolveArgs() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("resolve cmd args err invalid input: %v", cmd.input))
		return
	}
	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if suppliedArg := args.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.apiOptions[suppliedArg.Name()] = suppliedArg.Value()
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", argKV))
			}
		}
	}
}
