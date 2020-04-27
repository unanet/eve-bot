package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

type EvebotHelpCommand struct {
	baseCommand
}

func NewEvebotHelpCommand() EvebotHelpCommand {
	return EvebotHelpCommand{baseCommand{
		name:    "help",
		summary: "Try running one of the commands below",
		usage: bothelp.HelpUsage{
			"{{ command }} help",
		},
		optionalArgs:  botargs.Args{},
		examples:      bothelp.HelpExamples{},
		asyncRequired: false,
	}}
}

func (cmd EvebotHelpCommand) AsyncRequired() bool {
	return cmd.asyncRequired
}

func (cmd EvebotHelpCommand) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd EvebotHelpCommand) Name() string {
	return cmd.name
}

func (cmd EvebotHelpCommand) Help() *bothelp.Help {

	var nonHelpCmds string
	var nonHelpCmdExamples = bothelp.HelpExamples{}

	for _, v := range EvebotCommands {
		if v.Name() != cmd.name {
			nonHelpCmds = nonHelpCmds + "\n" + v.Name()
			nonHelpCmdExamples = append(nonHelpCmdExamples, v.Name()+" help")
		}
	}

	return bothelp.NewEvebotCommandHelp(
		bothelp.EvebotCommandHelpSummaryOpt(cmd.summary.String()),
		bothelp.EvebotCommandHelpCommandsOpt(nonHelpCmds),
		bothelp.EvebotCommandHelpUsageOpt(cmd.usage.String()),
		bothelp.EvebotCommandHelpArgsOpt(cmd.optionalArgs.String()),
		bothelp.EvebotCommandHelpExamplesOpt(nonHelpCmdExamples.String()),
	)

}

func (cmd EvebotHelpCommand) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd EvebotHelpCommand) AdditionalArgs() (botargs.Args, error) {
	return nil, nil
}

func (cmd EvebotHelpCommand) ResolveAdditionalArg(argKV []string) botargs.Arg {
	return nil
}
