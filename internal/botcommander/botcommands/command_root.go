package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

type RootCmd struct {
	baseCommand
}

func NewRootCmd(cmdFields []string, channel, user string) EvebotCommand {
	return RootCmd{baseCommand{
		input:          cmdFields,
		channel:        channel,
		user:           user,
		name:           "",
		summary:        "Welcome to `@evebot`! To get started, run:\n```@evebot help```",
		usage:          bothelp.Usage{},
		examples:       bothelp.Examples{},
		async:          false,
		optionalArgs:   botargs.Args{},
		requiredParams: botparams.Params{},
	}}
}

func (cmd RootCmd) User() string {
	return cmd.user
}

func (cmd RootCmd) Channel() string {
	return cmd.channel
}

func (cmd RootCmd) EveReqObj(user string) interface{} {
	return nil
}

func (cmd RootCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd RootCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd RootCmd) IsValid() bool {
	return true
}

func (cmd RootCmd) MakeAsyncReq() bool {
	return false
}

func (cmd RootCmd) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd RootCmd) Name() string {
	return cmd.name
}

func (cmd RootCmd) Help() *bothelp.Help {
	return bothelp.New(bothelp.HeaderOpt(cmd.summary.String()))
}

func (cmd RootCmd) IsHelpRequest() bool {
	return true
}
