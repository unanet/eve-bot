package botcommands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

func NewInvalidCommand(cmdFields []string) InvalidCmd {
	cmd := DefaultInvalidCommand()
	cmd.input = cmdFields
	cmd.summary = bothelp.Summary(fmt.Sprintf("I don't know how to execute the `%s` command.\n\nTry running: ```@evebot help```\n", cmdFields))
	return cmd
}

type InvalidCmd struct {
	baseCommand
}

func DefaultInvalidCommand() InvalidCmd {
	return InvalidCmd{baseCommand{
		name:           "",
		summary:        "Not sure what to do...",
		usage:          bothelp.Usage{},
		examples:       bothelp.Examples{},
		async:          false,
		optionalArgs:   botargs.Args{},
		suppliedArgs:   botargs.Args{},
		requiredParams: botparams.Params{},
		suppliedParams: botparams.Params{},
	}}
}

func (cmd InvalidCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd InvalidCmd) AckMsg(userID string) string {
	return baseAckMsg(cmd, userID, cmd.input)
}

func (cmd InvalidCmd) MakeAsyncReq() bool {
	return false
}

func (cmd InvalidCmd) IsValid() bool {
	return false
}

func (cmd InvalidCmd) Name() string {
	return cmd.name
}

func (cmd InvalidCmd) Help() *bothelp.Help {

	var nonHelpCmds string
	var nonHelpCmdExamples = bothelp.Examples{}

	for _, v := range EvebotCommands {
		if v.Name() != "help" {
			nonHelpCmds = nonHelpCmds + "\n" + v.Name()
			nonHelpCmdExamples = append(nonHelpCmdExamples, v.Name()+" help")
		}
	}

	return bothelp.New(
		bothelp.HeaderOpt(cmd.summary.String()),
		bothelp.CommandsOpt(nonHelpCmds),
		bothelp.ExamplesOpt(nonHelpCmdExamples.String()),
	)

}

func (cmd InvalidCmd) IsHelpRequest() bool {
	return true
}
