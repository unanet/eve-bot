package commands

import (
	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
)

type restartCmd struct {
	baseCommand
}

const (
	RestartCmdName = "restart"
)

var (
	restartCmdHelpSummary = help.Summary("The `restart` command is used to restart a service in a namespace")
	restartCmdHelpUsage   = help.Usage{"restart {{ service }} in {{ namespace }} {{ environment }}"}
	restartCmdHelpExample = help.Examples{"restart api in current int"}
)

// NewRestartCommand creates a New RestartCmd that implements the EvebotCommand interface
func NewRestartCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := restartCmd{baseCommand{
		input: cmdFields,
		info: ChatInfo{
			User:          user,
			Channel:       channel,
			CommandName:   RestartCmdName,
			IsHelpRequest: isHelpCmd(cmdFields, RestartCmdName),
		},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 5, Max: 5},
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

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd restartCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd restartCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *restartCmd) resolveDynamicOptions() {
	cmd.verifyInput()
	if len(cmd.errs) > 0 {
		return
	}

	cmd.opts[params.ServiceName] = cmd.input[1]
	cmd.opts[params.NamespaceName] = cmd.input[3]
	cmd.opts[params.EnvironmentName] = cmd.input[4]
}
