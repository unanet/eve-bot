package commands

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/botcommander/resources"
)

type showCmd struct {
	baseCommand
}

const (
	// ShowCmdName id/key
	ShowCmdName = "show"
)

var (
	showCmdHelpSummary = help.Summary("The `show` command is used to show resources (environments,namespaces,services,metadata,jobs)")
	showCmdHelpUsage   = help.Usage{
		"show {{ resources }}",
		"show namespaces in {{ environment }}",
		"show services in {{ namespace }} {{ environment }}",
		"show metadata for {{ service }} in {{ namespace }} {{ environment }}",
		"show jobs in {{ namespace }} {{ environment }}",
	}
	showCmdHelpExample = help.Examples{
		"show environments",
		"show namespaces in int",
		"show services in current int",
		"show metadata for billing in current int",
		"show jobs in current int",
	}
)

// NewShowCommand creates a New ShowCmd that implements the EvebotCommand interface
func NewShowCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := showCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: ShowCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 2, Max: 7},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd showCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(showCmdHelpSummary.String()),
		help.UsageOpt(showCmdHelpUsage.String()),
		help.ExamplesOpt(showCmdHelpExample.String()),
	).String())
}

func (cmd showCmd) IsAuthenticated(chatUserFn chatUserInfoFn, db *dynamodb.DynamoDB) bool {
	return false
}


// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd showCmd) IsAuthorized(allowedChannel map[string]interface{}, chatChanFn chatChannelInfoFn, chatUserFn chatUserInfoFn, db *dynamodb.DynamoDB) bool {
	if validUserRoleCheck(ShowCmdName, cmd, chatUserFn, db){
		return true
	}

	return true
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd showCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd showCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *showCmd) resolveDynamicOptions() {
	cmd.verifyInput()
	cmd.initializeResource()
	if len(cmd.errs) > 0 {
		return
	}

	switch cmd.opts["resource"] {
	case resources.JobName, "jobs":
		// show jobs in {{namespace}} {{environment}}
		if len(cmd.input) != 5 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid show jobs: %v", cmd.input))
			return
		}
		cmd.opts[params.NamespaceName] = cmd.input[3]
		cmd.opts[params.EnvironmentName] = cmd.input[4]
		return
	case resources.EnvironmentName:
		// show environments
		if len(cmd.input) != 2 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid show environment: %v", cmd.input))
			return
		}
		//...doesn't have any additional requirements
		return
	case resources.NamespaceName:
		// show namespaces in {{environment}}
		if len(cmd.input) != 4 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid show namespace: %v", cmd.input))
			return
		}
		cmd.opts[params.EnvironmentName] = cmd.input[3]
		return
	case resources.ServiceName:
		// show services in {{namespace}} {{environment}}
		if len(cmd.input) != 5 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid show service: %v", cmd.input))
			return
		}
		cmd.opts[params.NamespaceName] = cmd.input[3]
		cmd.opts[params.EnvironmentName] = cmd.input[4]
		return
	case resources.MetadataName:
		// show metadata for {{ service }} in {{ namespace }} {{ environment }}
		if len(cmd.input) != 7 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid show metadata: %v", cmd.input))
			return
		}
		cmd.opts[params.ServiceName] = cmd.input[3]
		cmd.opts[params.NamespaceName] = cmd.input[5]
		cmd.opts[params.EnvironmentName] = cmd.input[6]
		return
	default:
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource supplied: %v", cmd.opts["resource"]))
		return
	}
}
