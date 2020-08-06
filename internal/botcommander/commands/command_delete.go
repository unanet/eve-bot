package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
)

type DeleteCmd struct {
	baseCommand
}

func NewDeleteCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := DeleteCmd{baseCommand{
		input:       cmdFields,
		chatDetails: ChatInfo{User: user, Channel: channel},
		name:        "delete",
		summary:     "The `delete` command is used to delete resource values (metadata)",
		usage: help.Usage{
			"delete {{ resources }} for {{ service }} in {{ namespace }} {{ environment }}",
		},
		examples: help.Examples{
			"delete metadata for unaneta in current una-int key",
			"delete metadata for unaneta in current una-int key key2 key3 keyN",
			"delete version for unaneta in current una-int",
		},
		apiOptions:  make(CommandOptions),
		inputBounds: InputLengthBounds{Min: 7, Max: -1},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

func (cmd DeleteCmd) Details() CommandDetails {
	return CommandDetails{
		Name:          cmd.name,
		IsValid:       cmd.ValidInputLength(),
		IsHelpRequest: isHelpRequest(cmd.input, cmd.name),
		AckMsgFn:      baseAckMsg(cmd, cmd.input),
		ErrMsgFn:      cmd.BaseErrMsg(),
	}
}

func (cmd DeleteCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return validChannelAuthCheck(cmd.chatDetails.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.apiOptions)
}

func (cmd DeleteCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd DeleteCmd) ChatInfo() ChatInfo {
	return cmd.chatDetails
}

func (cmd DeleteCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd *DeleteCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete command: %v", cmd.input))
		return
	}

	if resources.IsValidDelete(cmd.input[1]) {
		cmd.apiOptions["resource"] = cmd.input[1]
	} else {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete resource: %v", cmd.input))
	}

	if cmd.apiOptions["resource"] == nil {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource: %v", cmd.input))
	}

	if len(cmd.errs) > 0 {
		return
	}

	switch cmd.apiOptions["resource"] {
	case resources.MetadataName:
		// delete metadata for unaneta in current una-int key,key2,key3
		// delete metadata for {{ service }} in {{ namespace }} {{ environment }} key,key2,key3
		if cmd.ValidInputLength() == false {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete metadata: %v", cmd.input))
			return
		}
		cmd.apiOptions[params.ServiceName] = cmd.input[3]
		cmd.apiOptions[params.NamespaceName] = cmd.input[5]
		cmd.apiOptions[params.EnvironmentName] = cmd.input[6]
		cmd.apiOptions[params.MetadataName] = cmd.input[7:]
		return
	case resources.VersionName:
		// delete version for unaneta in current una-int
		if cmd.ValidInputLength() == false {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid delete version: %v", cmd.input))
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
