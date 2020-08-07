package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type InvalidCmd struct {
	baseCommand
}

func NewInvalidCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := InvalidCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: ""},
		arguments:  args.Args{},
		parameters: params.Params{},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 1, Max: -1},
	}}
	return cmd
}

func (cmd InvalidCmd) AckMsg() (string, bool) {
	summary := help.Summary(fmt.Sprintf("I don't know how to execute the `%s` command.\n\nTry running: ```@evebot help```\n", cmd.input))

	helpMsg := help.New(
		help.HeaderOpt(summary.String()),
		help.CommandsOpt(NonHelpCmds),
		help.ExamplesOpt(NonHelpCommandExamples.String()),
	).String()

	return cmd.BaseAckMsg(helpMsg)
}

func (cmd InvalidCmd) IsAuthorized(map[string]interface{}, chatChannelInfoFn) bool {
	return true
}

func (cmd InvalidCmd) DynamicOptions() CommandOptions {
	return cmd.opts
}

func (cmd InvalidCmd) ChatInfo() ChatInfo {
	return cmd.info
}
