package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
)

type restartCmd struct {
	baseCommand
}

const (
	RestartCmdName = "restart"
)

var (
	restartCmdHelpSummary = help.Summary("The `restart` command is used to restart services in a namespace")
	restartCmdHelpUsage   = help.Usage{
		"restart {{ namespace }} in {{ environment }}",
		"restart {{ namespace }} in {{ environment }} services={{ service_name }}",
		"restart {{ namespace }} in {{ environment }} services={{ service_name,service_name,... }}",
	}
	restartCmdHelpExample = help.Examples{
		"restart current in una-int",
		"restart current in una-int services=subcontractor",
		"restart current in una-int services=subcontractor,platform",
	}
)

// NewRestartCommand creates a New RestartCmd that implements the EvebotCommand interface
func NewRestartCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := restartCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: RestartCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 4, Max: 5},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd restartCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(restartCmdHelpSummary.String()),
		help.UsageOpt(restartCmdHelpUsage.String()),
		help.ExamplesOpt(restartCmdHelpExample.String()),
	).String())
}

// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd restartCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return validChannelAuthCheck(cmd.info.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.opts)
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd restartCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd restartCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *restartCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid restart command: %v", cmd.input))
		return
	}

	if len(cmd.errs) > 0 {
		return
	}

	cmd.opts[params.NamespaceName] = cmd.input[1]
	cmd.opts[params.EnvironmentName] = cmd.input[3]

	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if suppliedArg := args.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.opts[suppliedArg.Name()] = suppliedArg.Value()
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", argKV))
			}
		}
	}
}
