package commander

type EvebotHelpCommand struct {
	name         string
	summary      EvebotCommandSummary
	usage        EvebotCommandUsage
	optionalArgs EvebotCommandArgs
	examples     EvebotCommandExamples
	input        []string
}

func NewEvebotHelpCommand() EvebotHelpCommand {
	return EvebotHelpCommand{
		name:    "help",
		summary: "Try running one of the commands below",
		usage: EvebotCommandUsage{
			"{{ command }} help",
		},
		optionalArgs: EvebotCommandArgs{},
		examples:     EvebotCommandExamples{},
	}
}

func (ehc EvebotHelpCommand) Initialize(input []string) EvebotCommand {
	ehc.input = input

	return ehc
}

func (ehc EvebotHelpCommand) Name() string {
	return ehc.name
}

func (ehc EvebotHelpCommand) Help() *EvebotCommandHelp {

	var nonHelpCmds string
	var nonHelpCmdExamples string

	for _, cmd := range evebotCommands {
		if cmd.Name() != ehc.name {
			nonHelpCmds = nonHelpCmds + "\n" + cmd.Name()
			nonHelpCmdExamples = nonHelpCmdExamples + "\n" + cmd.Name() + " help"
		}
	}

	return NewEvebotCommandHelp(
		EvebotCommandHelpSummaryOpt(ehc.summary.String()),
		EvebotCommandHelpCommandsOpt(nonHelpCmds),
		EvebotCommandHelpUsageOpt(ehc.usage.String()),
		EvebotCommandHelpArgsOpt(ehc.optionalArgs.String()),
		EvebotCommandHelpExamplesOpt(nonHelpCmdExamples),
	)

}

func (ehc EvebotHelpCommand) IsHelpRequest() bool {
	if len(ehc.input) == 0 || ehc.input[0] == "help" {
		return true
	}
	return false
}

func (ehc EvebotHelpCommand) AdditionalArgs() (EvebotCommandArgs, error) {
	return nil, nil
}

func (ehc EvebotHelpCommand) ResolveAdditionalArg(argKV []string) EvebotCommandArg {
	return nil
}
