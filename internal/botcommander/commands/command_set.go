package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
)

type SetCmd struct {
	baseCommand
}

func NewSetCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := SetCmd{baseCommand{
		input:       cmdFields,
		chatDetails: ChatInfo{User: user, Channel: channel},
		name:        "set",
		summary:     "The `set` command is used to set resource values (metadata and version)",
		usage: help.Usage{
			"set {{ resources }} for {{ service }} in {{ namespace }} {{ environment }} {{key=value}}",
			"set {{ resources }} in {{ namespace }} {{ environment }} to {{value}}",
		},
		examples: help.Examples{
			"set metadata for unaneta in current una-int key=value",
			"set metadata for unaneta in current una-int key=value key2=value2 keyN=valueN",
			"set version for unaneta in current una-int to 20.2",
			"set version in current una-int to 20.2",
		},
		apiOptions:  make(CommandOptions),
		inputBounds: InputLengthBounds{Min: 7, Max: -1},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

func (cmd SetCmd) Details() CommandDetails {
	return CommandDetails{
		Name:          cmd.name,
		IsValid:       cmd.ValidInputLength(),
		IsHelpRequest: isHelpRequest(cmd.input, cmd.name),
		AckMsgFn:      baseAckMsg(cmd, cmd.input),
		ErrMsgFn:      cmd.BaseErrMsg(),
	}
}

func (cmd SetCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return validChannelAuthCheck(cmd.chatDetails.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.apiOptions)
}

func (cmd SetCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd SetCmd) ChatInfo() ChatInfo {
	return cmd.chatDetails
}

func (cmd SetCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd *SetCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid set command params: %v", cmd.input))
		return
	}

	if resources.IsValidSet(cmd.input[1]) {
		cmd.apiOptions["resource"] = cmd.input[1]
	} else {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid set resource: %v", cmd.input))
		return
	}

	if cmd.apiOptions["resource"] == nil {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource: %v", cmd.input))
		return
	}

	if len(cmd.errs) > 0 {
		return
	}

	switch cmd.apiOptions["resource"] {
	case resources.MetadataName:
		// set metadata for unaneta in current una-int key=value
		// set metadata for {{ service }} in {{ namespace }} {{ environment }} key=value key=value key=value
		if len(cmd.input) < 8 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid set metadata: %v", cmd.input))
			return
		}
		cmd.apiOptions[params.ServiceName] = cmd.input[3]
		cmd.apiOptions[params.NamespaceName] = cmd.input[5]
		cmd.apiOptions[params.EnvironmentName] = cmd.input[6]
		cmd.apiOptions[params.MetadataName] = hydrateMetadataMap(cmd.input[7:])
		return
	case resources.VersionName:
		switch len(cmd.input) {
		// set version for unaneta in current una-int to 20.2
		case 9:
			cmd.apiOptions[params.ServiceName] = cmd.input[3]
			cmd.apiOptions[params.NamespaceName] = cmd.input[5]
			cmd.apiOptions[params.EnvironmentName] = cmd.input[6]
			cmd.apiOptions[params.VersionName] = cmd.input[8]
		// set version in current una-int to 20.2
		case 7:
			cmd.apiOptions[params.NamespaceName] = cmd.input[3]
			cmd.apiOptions[params.EnvironmentName] = cmd.input[4]
			cmd.apiOptions[params.VersionName] = cmd.input[6]
		default:
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid set version: %v", cmd.input))
			return
		}

	default:
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource supplied: %v", cmd.apiOptions["resource"]))
		return
	}
}
