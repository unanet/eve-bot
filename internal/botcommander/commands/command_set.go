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

const (
	SetCmdName = "set"
)

var (
	setCmdHelpSummary = help.Summary("The `set` command is used to set resource values (metadata and version)")
	setCmdHelpUsage   = help.Usage{
		"set {{ resources }} for {{ service }} in {{ namespace }} {{ environment }} {{key=value}}",
		"set {{ resources }} in {{ namespace }} {{ environment }} to {{value}}",
	}
	setCmdHelpExample = help.Examples{
		"set metadata for unaneta in current una-int key=value",
		"set metadata for unaneta in current una-int key=value key2=value2 keyN=valueN",
		"set version for unaneta in current una-int to 20.2",
		"set version in current una-int to 20.2",
	}
)

func NewSetCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := SetCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: SetCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 7, Max: -1},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

func (cmd SetCmd) AckMsg() (string, bool) {

	helpMsg := help.New(
		help.HeaderOpt(setCmdHelpSummary.String()),
		help.UsageOpt(setCmdHelpUsage.String()),
		help.ExamplesOpt(setCmdHelpExample.String()),
	).String()

	return cmd.BaseAckMsg(helpMsg)
}

func (cmd SetCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return validChannelAuthCheck(cmd.info.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.opts)
}

func (cmd SetCmd) DynamicOptions() CommandOptions {
	return cmd.opts
}

func (cmd SetCmd) ChatInfo() ChatInfo {
	return cmd.info
}

func (cmd *SetCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid set command params: %v", cmd.input))
		return
	}

	if resources.IsValidSet(cmd.input[1]) {
		cmd.opts["resource"] = cmd.input[1]
	} else {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid set resource: %v", cmd.input))
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
	case resources.MetadataName:
		// set metadata for unaneta in current una-int key=value
		// set metadata for {{ service }} in {{ namespace }} {{ environment }} key=value key=value key=value
		if len(cmd.input) < 8 {
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid set metadata: %v", cmd.input))
			return
		}
		cmd.opts[params.ServiceName] = cmd.input[3]
		cmd.opts[params.NamespaceName] = cmd.input[5]
		cmd.opts[params.EnvironmentName] = cmd.input[6]
		cmd.opts[params.MetadataName] = hydrateMetadataMap(cmd.input[7:])
		return
	case resources.VersionName:
		switch len(cmd.input) {
		// set version for unaneta in current una-int to 20.2
		case 9:
			cmd.opts[params.ServiceName] = cmd.input[3]
			cmd.opts[params.NamespaceName] = cmd.input[5]
			cmd.opts[params.EnvironmentName] = cmd.input[6]
			cmd.opts[params.VersionName] = cmd.input[8]
		// set version in current una-int to 20.2
		case 7:
			cmd.opts[params.NamespaceName] = cmd.input[3]
			cmd.opts[params.EnvironmentName] = cmd.input[4]
			cmd.opts[params.VersionName] = cmd.input[6]
		default:
			cmd.errs = append(cmd.errs, fmt.Errorf("invalid set version: %v", cmd.input))
			return
		}

	default:
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid resource supplied: %v", cmd.opts["resource"]))
		return
	}
}
