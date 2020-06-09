package botcommands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

func NewInvalidCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultInvalidCommand(cmdFields, channel, user)
}

type InvalidCmd struct {
	baseCommand
}

func defaultInvalidCommand(cmdFields []string, channel, user string) InvalidCmd {
	return InvalidCmd{baseCommand{
		input:          cmdFields,
		channel:        channel,
		user:           user,
		name:           "",
		summary:        bothelp.Summary(fmt.Sprintf("I don't know how to execute the `%s` command.\n\nTry running: ```@evebot help```\n", cmdFields)),
		usage:          bothelp.Usage{},
		examples:       bothelp.Examples{},
		async:          false,
		optionalArgs:   botargs.Args{},
		requiredParams: botparams.Params{},
	}}
}

func (cmd InvalidCmd) User() string {
	return cmd.user
}

func (cmd InvalidCmd) Channel() string {
	return cmd.channel
}

func (cmd InvalidCmd) EveReqObj(user string) interface{} {
	return nil
}

func (cmd InvalidCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd InvalidCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
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

	for _, v := range nonHelpCmd() {
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
