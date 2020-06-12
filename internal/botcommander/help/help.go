package help

type Help struct {
	Summary  string
	Usage    string
	Examples string
	Args     string
	Commands string
	Header   string
}

// Evebot Command Help
func (h Help) String() string {
	var msg string
	if len(h.Header) > 0 {
		msg = msg + "\n" + h.Header + "\n\n"
	}
	if len(h.Summary) > 0 {
		msg = msg + "*Summary:* " + h.Summary + "...\n\n"
	}
	if len(h.Commands) > 0 {
		msg = msg + "*Commands:*\n" + "```" + h.Commands + "```" + "\n\n"
	}
	if len(h.Usage) > 0 {
		msg = msg + "*Usage:*\n" + "```" + h.Usage + "```" + "\n\n"
	}
	if len(h.Args) > 0 {
		msg = msg + "*Optional Args:*\n" + "```" + h.Args + "```" + "\n\n"
	}
	if len(h.Examples) > 0 {
		msg = msg + "*Examples:*\n" + "```" + h.Examples + "```" + "\n\n"
	}
	return msg
}

type HelpOption func(*Help)

func New(opts ...HelpOption) *Help {
	h := &Help{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func HeaderOpt(header string) HelpOption {
	return func(h *Help) {
		h.Header = header
	}
}

func SummaryOpt(summary string) HelpOption {
	return func(h *Help) {
		h.Summary = summary
	}
}

func UsageOpt(usage string) HelpOption {
	return func(h *Help) {
		h.Usage = usage
	}
}

func ArgsOpt(args string) HelpOption {
	return func(h *Help) {
		h.Args = args
	}
}

func ExamplesOpt(examples string) HelpOption {
	return func(h *Help) {
		h.Examples = examples
	}
}

func CommandsOpt(commands string) HelpOption {
	return func(h *Help) {
		h.Commands = commands
	}
}
