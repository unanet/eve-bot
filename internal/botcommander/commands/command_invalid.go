package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type invalidCmd struct {
	baseCommand
}

// NewInvalidCommand creates a New InvalidCmd that implements the EvebotCommand interface
func NewInvalidCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := invalidCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: ""},
		arguments:  args.Args{},
		parameters: params.Params{},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 1, Max: -1},
	}}
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd invalidCmd) AckMsg() (string, bool) {
	summary := help.Summary(fmt.Sprintf("I don't know how to execute the `%s` command.\n\nTry running: ```@evebot help```\n", cmd.input)).String()
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(summary),
		help.CommandsOpt(NonHelpCmds),
		help.ExamplesOpt(NonHelpCommandExamples.String()),
	).String())
}

// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd invalidCmd) IsAuthorized(map[string]interface{}, chatChannelInfoFn) bool {
	return true
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd invalidCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd invalidCmd) Info() ChatInfo {
	return cmd.info
}
