package commands

import (
	"fmt"

	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/botcommander/resources"
)

type setCmd struct {
	baseCommand
}

const (
	// SetCmdName id/key
	SetCmdName = "set"
)

var (
	setCmdHelpSummary = help.Summary("The `set` command is used to set resource values (metadata and version)")
	setCmdHelpUsage   = help.Usage{
		"set {{ resources }} for {{ service }} in {{ namespace }} {{ environment }} {{key=value}}",
		"set {{ resources }} in {{ namespace }} {{ environment }} to {{value}}",
	}
	setCmdHelpExample = help.Examples{
		"set metadata for api in current int key=value",
		"set metadata for billing in current int key=value key2=value2 keyN=valueN",
		"set version for api in current int to 1.3",
		"set version in current int to 2.0",
	}
)

// NewSetCommand creates a New SetCmd that implements the EvebotCommand interface
func NewSetCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := setCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: SetCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 7, Max: -1},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd setCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(setCmdHelpSummary.String()),
		help.UsageOpt(setCmdHelpUsage.String()),
		help.ExamplesOpt(setCmdHelpExample.String()),
	).String())
}

// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd setCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return cmd.IsHelpRequest() || validChannelAuthCheck(cmd.info.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.opts)
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd setCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd setCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *setCmd) resolveDynamicOptions() {
	cmd.verifyInput()
	cmd.initializeResource()
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
