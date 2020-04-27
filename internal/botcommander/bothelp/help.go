package bothelp

type Help struct {
	Summary  string
	Usage    string
	Examples string
	Args     string
	Commands string
	Header   string
}

// Evebot Command Help
func (ebch Help) String() string {
	var msg string
	if len(ebch.Header) > 0 {
		msg = msg + "\n" + ebch.Header + "\n\n"
	}
	if len(ebch.Summary) > 0 {
		msg = msg + "*Summary:* " + ebch.Summary + "...\n\n"
	}
	if len(ebch.Commands) > 0 {
		msg = msg + "*Commands:*\n" + "```" + ebch.Commands + "```" + "\n\n"
	}
	if len(ebch.Usage) > 0 {
		msg = msg + "*Usage:*\n" + "```" + ebch.Usage + "```" + "\n\n"
	}
	if len(ebch.Args) > 0 {
		msg = msg + "*Optional Args:*\n" + "```" + ebch.Args + "```" + "\n\n"
	}
	if len(ebch.Examples) > 0 {
		msg = msg + "*Examples:*\n" + "```" + ebch.Examples + "```" + "\n\n"
	}
	return msg
}

type HelpOption func(ech *Help)

func EvebotCommandHelpHeaderOpt(header string) HelpOption {
	return func(ech *Help) {
		ech.Header = header
	}
}

func NewEvebotCommandHelp(opts ...HelpOption) *Help {
	e := &Help{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func EvebotCommandHelpSummaryOpt(summary string) HelpOption {
	return func(ech *Help) {
		ech.Summary = summary
	}
}

func EvebotCommandHelpUsageOpt(usage string) HelpOption {
	return func(ech *Help) {
		ech.Usage = usage
	}
}

func EvebotCommandHelpArgsOpt(args string) HelpOption {
	return func(ech *Help) {
		ech.Args = args
	}
}

func EvebotCommandHelpExamplesOpt(examples string) HelpOption {
	return func(ech *Help) {
		ech.Examples = examples
	}
}

func EvebotCommandHelpCommandsOpt(commands string) HelpOption {
	return func(ech *Help) {
		ech.Commands = commands
	}
}
