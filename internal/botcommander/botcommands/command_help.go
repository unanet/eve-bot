package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

func NewHelpCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultHelpCommand(cmdFields, channel, user)
}

func (cmd HelpCmd) EveReqObj(user string) interface{} {
	return nil
}

type HelpCmd struct {
	baseCommand
}

func defaultHelpCommand(cmdFields []string, channel, user string) HelpCmd {
	return HelpCmd{baseCommand{
		input:   cmdFields,
		channel: channel,
		user:    user,
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

func (cmd HelpCmd) User() string {
	return cmd.user
}

func (cmd HelpCmd) Channel() string {
	return cmd.channel
}

func (cmd HelpCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
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
