package commander

type EvebotHelpCommand struct {
	name         string
	summary      EvebotCommandSummary
	usage        EvebotCommandUsage
	optionalArgs EvebotCommandArgs
	examples     EvebotCommandExamples
	commandList  EvebotCommandList
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

func (ehc EvebotHelpCommand) Examples() EvebotCommandExamples {
	return ehc.examples
}

func (ehc EvebotHelpCommand) Name() string {
	return ehc.name
}

func (ehc EvebotHelpCommand) OptionalArgs() EvebotCommandArgs {
	return ehc.optionalArgs
}

func (ehc EvebotHelpCommand) Summary() EvebotCommandSummary {
	return ehc.summary
}

func (ehc EvebotHelpCommand) Usage() EvebotCommandUsage {
	return ehc.usage
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

func (ehc EvebotHelpCommand) IsHelpRequest(input []string) bool {
	if len(input) == 0 || input[0] == "help" {
		return true
	}
	return false
}

func (ehc EvebotHelpCommand) IsValidCommand(input []string) bool {
	if len(input) <= 0 || input[0] != ehc.Name() || len(input) >= 3 {
		return false
	}
	return true
}

func (ehc EvebotHelpCommand) AdditionalArgs(input []string) (EvebotCommandArgs, error) {
	return nil, nil
}

func (ehc EvebotHelpCommand) ResolveAdditionalArg(argKV []string) EvebotArg {
	return nil
}
