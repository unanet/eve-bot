package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
)

func NewDeleteCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultDeleteCommand(cmdFields, channel, user)
}

type DeleteCmd struct {
	baseCommand
}

func defaultDeleteCommand(cmdFields []string, channel, user string) DeleteCmd {
	cmd := DeleteCmd{baseCommand{
		input:   cmdFields,
		channel: channel,
		user:    user,
		name:    "delete",
		summary: "The `delete` command is used to delete resource values (metadata)",
		usage: help.Usage{
			"delete {{ resources }} for {{ service }} in {{ namespace }} {{ environment }}",
		},
		examples: help.Examples{
			"delete metadata for unaneta in current una-int key",
			"delete metadata for unaneta in current una-int key,key2,key3",
		},
		apiOptions:          make(CommandOptions),
		requiredInputLength: 7,
	}}
	cmd.resolveResource()
	cmd.resolveConditionalParams()
	return cmd
}

func (cmd DeleteCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd DeleteCmd) User() string {
	return cmd.user
}

func (cmd DeleteCmd) Channel() string {
	return cmd.channel
}

func (cmd DeleteCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd DeleteCmd) IsValid() bool {
	if baseIsValid(cmd.input) && len(cmd.input) >= cmd.requiredInputLength {
		return true
	}
	return false
}

func (cmd DeleteCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd DeleteCmd) Name() string {
	return cmd.name
}

func (cmd DeleteCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd DeleteCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd *DeleteCmd) resolveResource() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete command: %v", cmd.input))
		return
	}

	if resources.IsValidDelete(cmd.input[1]) {
		cmd.apiOptions["resource"] = cmd.input[1]
	} else {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete resource: %v", cmd.input))
		return
	}

}

func (cmd *DeleteCmd) resolveConditionalParams() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete command params: %v", cmd.input))
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
	case resources.MetadataName:
		// delete metadata for unaneta in current una-int key,key2,key3
		// delete metadata for {{ service }} in {{ namespace }} {{ environment }} key,key2,key3
		if len(cmd.input) < 7 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete metadata: %v", cmd.input))
			return
		}
		cmd.apiOptions[params.ServiceName] = cmd.input[3]
		cmd.apiOptions[params.NamespaceName] = cmd.input[5]
		cmd.apiOptions[params.EnvironmentName] = cmd.input[6]
		cmd.apiOptions[params.MetadataName] = strings.Split(cmd.input[7:][0], ",")
		return
	default:
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource supplied: %v", cmd.apiOptions["resource"]))
		return
	}
}
