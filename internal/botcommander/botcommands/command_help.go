package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

type HelpCmd struct {
	baseCommand
}

func NewEvebotHelpCommand() HelpCmd {
	return HelpCmd{baseCommand{
		name:    "help",
		summary: "Try running one of the commands below",
		usage: bothelp.Usage{
			"{{ command }} help",
		},
		optionalArgs:  botargs.Args{},
		examples:      bothelp.Examples{},
		asyncRequired: false,
	}}
}

func (cmd HelpCmd) AsyncRequired() bool {
	return cmd.asyncRequired
}

func (cmd HelpCmd) IsValid() bool {
	if cmd.input == nil || len(cmd.input) == 0 {
		return false
	}
	return true
}

func (cmd HelpCmd) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd HelpCmd) Name() string {
	return cmd.name
}

func (cmd HelpCmd) Help() *bothelp.Help {

	var nonHelpCmds string
	var nonHelpCmdExamples = bothelp.Examples{}

	for _, v := range EvebotCommands {
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

func (cmd HelpCmd) AdditionalArgs() (botargs.Args, error) {
	return nil, nil
}

func (cmd HelpCmd) ResolveAdditionalArg(argKV []string) botargs.Arg {
	return nil
}
