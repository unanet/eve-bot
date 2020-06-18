package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
)

func NewShowCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultShowCommand(cmdFields, channel, user)
}

type ShowCmd struct {
	baseCommand
}

// @evebot show current in qa
func defaultShowCommand(cmdFields []string, channel, user string) ShowCmd {
	cmd := ShowCmd{baseCommand{
		input:   cmdFields,
		channel: channel,
		user:    user,
		name:    "show",
		summary: "The `show` command is used to show resources (environments,namespaces,services,metadata)",
		usage: help.Usage{
			"show {{ resources }}",
			"show namespaces in {{ environment }}",
			"show services in {{ namespace }} {{ environment }}",
			"show metadata for {{ service }} in {{ namespace }} {{ environment }}",
		},
		examples: help.Examples{
			"show environments",
			"show namespaces in una-int",
			"show services in current una-int",
			"show metadata for unaneta in current una-int",
		},
		apiOptions:          make(CommandOptions),
		requiredInputLength: 2,
	}}
	cmd.resolveResource()
	cmd.resolveConditionalParams()
	return cmd
}

func (cmd ShowCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd ShowCmd) User() string {
	return cmd.user
}

func (cmd ShowCmd) Channel() string {
	return cmd.channel
}

func (cmd ShowCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd ShowCmd) IsValid() bool {
	if baseIsValid(cmd.input) && len(cmd.input) >= cmd.requiredInputLength {
		return true
	}
	return false
}

func (cmd ShowCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd ShowCmd) Name() string {
	return cmd.name
}

func (cmd ShowCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd ShowCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd *ShowCmd) resolveResource() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid show command: %v", cmd.input))
		return
	}

	if resources.IsValid(cmd.input[1]) {
		cmd.apiOptions["resource"] = cmd.input[1]
	} else {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid requested resource: %v", cmd.input))
		return
	}

}

func (cmd *ShowCmd) resolveConditionalParams() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid show command params: %v", cmd.input))
		return
	}

	if len(cmd.errs) > 0 {
		return
	}

	if cmd.apiOptions["resource"] == nil {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource: %v", cmd.input))
		return
	}

	switch cmd.apiOptions["resource"] {
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
		cmd.apiOptions[params.EnvironmentName] = cmd.input[3]
		return
	case resources.ServiceName:
		// show services in {{namespace}} {{environment}}
		if len(cmd.input) != 5 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid show service: %v", cmd.input))
			return
		}
		cmd.apiOptions[params.NamespaceName] = cmd.input[3]
		cmd.apiOptions[params.EnvironmentName] = cmd.input[4]
		return
	case resources.MetadataName:
		// show metadata for {{ service }} in {{ namespace }} {{ environment }}
		if len(cmd.input) != 7 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid show metadata: %v", cmd.input))
			return
		}
		cmd.apiOptions[params.ServiceName] = cmd.input[3]
		cmd.apiOptions[params.NamespaceName] = cmd.input[5]
		cmd.apiOptions[params.EnvironmentName] = cmd.input[6]
		return
	default:
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource supplied: %v", cmd.apiOptions["resource"]))
		return
	}
}
