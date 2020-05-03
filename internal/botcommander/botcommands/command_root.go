package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

type RootCmd struct {
	baseCommand
}

func NewRootCmd() EvebotCommand {
	return RootCmd{baseCommand{
		name:           "",
		summary:        "Welcome to `@evebot`! To get started, run:\n```@evebot help```",
		usage:          bothelp.Usage{},
		examples:       bothelp.Examples{},
		async:          false,
		optionalArgs:   botargs.Args{},
		requiredParams: botparams.Params{},
	}}
}

func (cmd RootCmd) EveReqObj() interface{} {
	return nil
}

func (cmd RootCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd RootCmd) AckMsg(userID string) string {
	return baseAckMsg(cmd, userID, cmd.input)
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
	return bothelp.New(
		bothelp.HeaderOpt(cmd.summary.String()),
	)
}

func (cmd RootCmd) IsHelpRequest() bool {
	return true
}
