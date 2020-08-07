package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
)

type ShowCmd struct {
	baseCommand
}

const (
	ShowCmdName = "show"
)

var (
	showCmdHelpSummary = help.Summary("The `show` command is used to show resources (environments,namespaces,services,metadata)")

	showCmdHelpUsage = help.Usage{
		"show {{ resources }}",
		"show namespaces in {{ environment }}",
		"show services in {{ namespace }} {{ environment }}",
		"show metadata for {{ service }} in {{ namespace }} {{ environment }}",
	}
	showCmdHelpExample = help.Examples{
		"show environments",
		"show namespaces in una-int",
		"show services in current una-int",
		"show metadata for unaneta in current una-int",
	}
)

func NewShowCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := ShowCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: ShowCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 2, Max: 7},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

func (cmd ShowCmd) AckMsg() (string, bool) {

	helpMsg := help.New(
		help.HeaderOpt(showCmdHelpSummary.String()),
		help.UsageOpt(showCmdHelpUsage.String()),
		help.ExamplesOpt(showCmdHelpExample.String()),
	).String()

	return cmd.BaseAckMsg(helpMsg)
}

func (cmd ShowCmd) IsAuthorized(map[string]interface{}, chatChannelInfoFn) bool {
	return true
}

func (cmd ShowCmd) DynamicOptions() CommandOptions {
	return cmd.opts
}

func (cmd ShowCmd) ChatInfo() ChatInfo {
	return cmd.info
}

func (cmd *ShowCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid show command: %v", cmd.input))
		return
	}

	if resources.IsValid(cmd.input[1]) {
		cmd.opts["resource"] = cmd.input[1]
	} else {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid requested resource: %v", cmd.input))
		return
	}

	if cmd.opts["resource"] == nil {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource: %v", cmd.input))
		return
	}

	if len(cmd.errs) > 0 {
		return
	}

	switch cmd.opts["resource"] {
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
