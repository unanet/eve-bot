package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

func NewHelpCommand(cmdFields []string) EvebotCommand {
	cmd := defaultHelpCommand()
	cmd.input = cmdFields
	return cmd
}

func (cmd HelpCmd) EveReqObj(cbURL, user string) interface{} {
	return nil
}

type HelpCmd struct {
	baseCommand
}

func defaultHelpCommand() HelpCmd {
	return HelpCmd{baseCommand{
		name:    "help",
		summary: "Try running one of the commands below",
		usage: bothelp.Usage{
			"{{ command }} help",
		},
		examples:       bothelp.Examples{},
		async:          false,
		optionalArgs:   botargs.Args{},
		requiredParams: botparams.Params{},
	}}
}

func (cmd HelpCmd) AckMsg(userID string) string {
	return baseAckMsg(cmd, userID, cmd.input)
}

func (cmd HelpCmd) MakeAsyncReq() bool {
	return false
}

func (cmd HelpCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd HelpCmd) IsValid() bool {
	if len(cmd.errs) > 0 {
		return false
	}
	return baseIsValid(cmd.input)
}

func (cmd HelpCmd) Name() string {
	return cmd.name
}

func (cmd HelpCmd) Help() *bothelp.Help {
	var nonHelpCmds string
	var nonHelpCmdExamples = bothelp.Examples{}

	for _, v := range nonHelpCmd() {
		if v.Name() != cmd.name {
			nonHelpCmds = nonHelpCmds + "\n" + v.Name()
			nonHelpCmdExamples = append(nonHelpCmdExamples, v.Name()+" help")
		}
	}

	return bothelp.New(
		bothelp.HeaderOpt(cmd.summary.String()),
		bothelp.CommandsOpt(nonHelpCmds),
		bothelp.UsageOpt(cmd.usage.String()),
		bothelp.ArgsOpt(cmd.optionalArgs.String()),
		bothelp.ExamplesOpt(nonHelpCmdExamples.String()),
	)
}

func (cmd HelpCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}
