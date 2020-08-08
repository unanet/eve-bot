package help

// Help data structure
type Help struct {
	Summary  string
	Usage    string
	Examples string
	Args     string
	Commands string
	Header   string
}

// String Evebot Command Help
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

// Option type for dynamic opts
type Option func(*Help)

// New creates a new Help structure
func New(opts ...Option) *Help {
	h := &Help{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// HeaderOpt is a header option
func HeaderOpt(header string) Option {
	return func(h *Help) {
		h.Header = header
	}
}

// UsageOpt is a header option
func UsageOpt(usage string) Option {
	return func(h *Help) {
		h.Usage = usage
	}
}

// ArgsOpt is a header option
func ArgsOpt(args string) Option {
	return func(h *Help) {
		h.Args = args
	}
}

// ExamplesOpt is a header option
func ExamplesOpt(examples string) Option {
	return func(h *Help) {
		h.Examples = examples
	}
}

// CommandsOpt is a header option
func CommandsOpt(commands string) Option {
	return func(h *Help) {
		h.Commands = commands
	}
}
