package commander

type EvebotCommandHelp struct {
	Summary  string
	Usage    string
	Examples string
	Args     string
	Commands string
}

// Evebot Command Help
func (ebch EvebotCommandHelp) String() string {
	var msg string
	if len(ebch.Summary) > 0 {
		msg = msg + "`Summary:`\n" + ebch.Summary + "\n\n"
	}
	if len(ebch.Commands) > 0 {
		msg = msg + "`Commands:`\n" + ebch.Commands + "\n\n"
	}
	if len(ebch.Usage) > 0 {
		msg = msg + "`Usage:`\n" + ebch.Usage + "\n\n"
	}
	if len(ebch.Args) > 0 {
		msg = msg + "`Optional Args:`\n" + ebch.Args + "\n\n"
	}
	if len(ebch.Examples) > 0 {
		msg = msg + "`Examples:`\n" + ebch.Examples + "\n\n"
	}
	return msg
}

type EvebotCommandHelpOption func(ech *EvebotCommandHelp)

func NewEvebotCommandHelp(opts ...EvebotCommandHelpOption) *EvebotCommandHelp {
	e := &EvebotCommandHelp{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func EvebotCommandHelpSummaryOpt(summary string) EvebotCommandHelpOption {
	return func(ech *EvebotCommandHelp) {
		ech.Summary = summary
	}
}

func EvebotCommandHelpUsageOpt(usage string) EvebotCommandHelpOption {
	return func(ech *EvebotCommandHelp) {
		ech.Usage = usage
	}
}

func EvebotCommandHelpArgsOpt(args string) EvebotCommandHelpOption {
	return func(ech *EvebotCommandHelp) {
		ech.Args = args
	}
}

func EvebotCommandHelpExamplesOpt(examples string) EvebotCommandHelpOption {
	return func(ech *EvebotCommandHelp) {
		ech.Examples = examples
	}
}

func EvebotCommandHelpCommandsOpt(commands string) EvebotCommandHelpOption {
	return func(ech *EvebotCommandHelp) {
		ech.Commands = commands
	}
}
