package commands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type HelpCmd struct {
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

func NewHelpCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := HelpCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: helpCmdName},
		arguments:  args.Args{},
		parameters: params.Params{},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 1, Max: -1},
	}}
	return cmd
}

func (cmd HelpCmd) AckMsg() (string, bool) {

	helpMsg := help.New(
		help.HeaderOpt(helpCmdHelpSummary.String()),
		help.CommandsOpt(NonHelpCmds),
		help.UsageOpt(helpCmdHelpUsage.String()),
		help.ArgsOpt(cmd.arguments.String()),
		help.ExamplesOpt(NonHelpCommandExamples.String()),
	).String()

	return cmd.BaseAckMsg(helpMsg)
}

func (cmd HelpCmd) IsAuthorized(map[string]interface{}, chatChannelInfoFn) bool {
	return true
}

func (cmd HelpCmd) DynamicOptions() CommandOptions {
	return cmd.opts
}

func (cmd HelpCmd) ChatInfo() ChatInfo {
	return cmd.info
}
