package commands

import (
	"github.com/unanet/eve-bot/internal/botcommander/args"
	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
)

type helpCmd struct {
	baseCommand
}

const (
	helpCmdName = "help"
)

var (
	helpCmdHelpSummary = help.Summary("Try running one of the commands below")
	helpCmdHelpUsage   = help.Usage{
		"help",
		"{{ command }} help",
	}
)

// NewHelpCommand creates a New HelpCmd that implements the EvebotCommand interface
func NewHelpCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := helpCmd{baseCommand{
		input: cmdFields,
		info: ChatInfo{
			User:          user,
			Channel:       channel,
			CommandName:   helpCmdName,
			IsHelpRequest: true,
		},
		arguments:  args.Args{},
		parameters: params.Params{},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 1, Max: -1},
	}}
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd helpCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(helpCmdHelpSummary.String()),
		help.CommandsOpt(NewFactory().NonHelpCmds()),
		help.UsageOpt(helpCmdHelpUsage.String()),
		help.ArgsOpt(cmd.arguments.String()),
		help.ExamplesOpt(NewFactory().NonHelpExamples().String()),
	).String())
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd helpCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd helpCmd) Info() ChatInfo {
	return cmd.info
}
