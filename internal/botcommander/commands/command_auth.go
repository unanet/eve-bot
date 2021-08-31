package commands

import (
	"github.com/unanet/eve-bot/internal/botcommander/help"
)

type authCmd struct {
	baseCommand
}

const (
	// AuthCmdName used as key/id for the auth command
	AuthCmdName = "auth"
)

var (
	authCmdHelpSummary = help.Summary("The `auth` command is used to authenticate")
	authCmdHelpUsage   = help.Usage{
		"auth",
	}
	authCmdHelpExample = help.Examples{}
)

// NewAuthCommand creates a New AuthCmd that implements the EvebotCommand interface
func NewAuthCommand(cmdFields []string, channel, user string) EvebotCommand {
	return authCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: AuthCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 1, Max: 1},
	}}
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd authCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(authCmdHelpSummary.String()),
		help.UsageOpt(authCmdHelpUsage.String()),
		help.ExamplesOpt(authCmdHelpExample.String()),
	).String())
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd authCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd authCmd) Info() ChatInfo {
	return cmd.info
}
