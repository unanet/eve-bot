package commands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type RootCmd struct {
	baseCommand
}

const (
	rootCmdHelpSummary = help.Summary("Welcome to `@evebot`! To get started, run:\n```@evebot help```")
)

func NewRootCmd(cmdFields []string, channel, user string) EvebotCommand {
	cmd := RootCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: ""},
		arguments:  args.Args{},
		parameters: params.Params{},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 0, Max: -1},
	}}
	return cmd
}

func (cmd RootCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(help.HeaderOpt(rootCmdHelpSummary.String())).String())
}

func (cmd RootCmd) IsAuthorized(map[string]interface{}, chatChannelInfoFn) bool {
	return true
}

func (cmd RootCmd) DynamicOptions() CommandOptions {
	return cmd.opts
}

func (cmd RootCmd) ChatInfo() ChatInfo {
	return cmd.info
}
