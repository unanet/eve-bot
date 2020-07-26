package commands

import (
	"fmt"
	"regexp"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
)

func NewSetCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultSetCommand(cmdFields, channel, user)
}

type SetCmd struct {
	baseCommand
}

func defaultSetCommand(cmdFields []string, channel, user string) SetCmd {
	cmd := SetCmd{baseCommand{
		input:   cmdFields,
		channel: channel,
		user:    user,
		name:    "set",
		summary: "The `set` command is used to set resource values (metadata and version)",
		usage: help.Usage{
			"set {{ resources }} for {{ service }} in {{ namespace }} {{ environment }}",
			"set {{ resources }} in {{ namespace }} {{ environment }}",
		},
		examples: help.Examples{
			"set metadata for unaneta in current una-int key=value",
			"set metadata for unaneta in current una-int key=value key2=value2 keyN=valueN",
			"set version for unaneta in current una-int to 20.2",
			"set version in current una-int to 20.2",
		},
		apiOptions:          make(CommandOptions),
		requiredInputLength: 4,
	}}
	cmd.resolveResource()
	cmd.resolveConditionalParams()
	return cmd
}

func (cmd SetCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd SetCmd) User() string {
	return cmd.user
}

func (cmd SetCmd) Channel() string {
	return cmd.channel
}

func (cmd SetCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd SetCmd) IsValid() bool {
	if baseIsValid(cmd.input) && len(cmd.input) >= cmd.requiredInputLength {
		return true
	}
	return false
}

func (cmd SetCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd SetCmd) Name() string {
	return cmd.name
}

func (cmd SetCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd SetCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd *SetCmd) resolveResource() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid set command: %v", cmd.input))
		return
	}

	if resources.IsValidSet(cmd.input[1]) {
		cmd.apiOptions["resource"] = cmd.input[1]
	} else {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid set resource: %v", cmd.input))
		return
	}

}

func (cmd *SetCmd) resolveConditionalParams() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid set command params: %v", cmd.input))
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

func (cmd SetCmd) validMetadataMap(metadataMap params.MetadataMap) bool {
	invalidCharMatcher := regexp.MustCompile(`<http:\/\/|>|<|\/\/|\||https:\/\/`)
	result := true
	for k, v := range metadataMap {
		invalidCharKeyMatchIndexes := invalidCharMatcher.FindAllStringIndex(k, -1)
		if len(invalidCharKeyMatchIndexes) > 0 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid metadata key supplied: %v", k))
			result = false
		}

		if strVal, ok := v.(string); ok {
			invalidCharValueMatchIndexes := invalidCharMatcher.FindAllStringIndex(strVal, -1)
			if len(invalidCharValueMatchIndexes) > 0 {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid metadata value supplied: %v", strVal))
				result = false
			}
		}
	}
	return result
}
