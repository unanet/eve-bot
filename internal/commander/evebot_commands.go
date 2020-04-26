package commander

func init() {
	// Add all the Evebot commands here on init
	evebotCommands = []EvebotCommand{
		NewEvebotHelpCommand(),
		NewEvebotDeployCommand(),
	}

	// Set the Full Help Message Once on init
	CmdHelpMsgs = fullHelpMessage()
}

var (
	evebotCommands []EvebotCommand
	CmdHelpMsgs    string
)

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

// Evebot Command List
type EvebotCommandList []string

func (ebcl EvebotCommandList) String() string {
	var msg string
	for _, s := range ebcl {
		if len(msg) > 0 {
			msg = msg + "\n" + s
		} else {
			msg = s
		}
	}
	return msg
}

// Evebot Command Examples
type EvebotCommandExamples []string

func (ebce EvebotCommandExamples) String() string {
	var msg string
	for _, s := range ebce {
		if len(msg) > 0 {
			msg = msg + "\n" + s
		} else {
			msg = s
		}
	}
	return msg
}

// Evebot Command Summary
type EvebotCommandSummary string

func (ebcs EvebotCommandSummary) String() string {
	return string(ebcs)
}

// Evebot Command Usage
type EvebotCommandUsage []string

func (ebcu EvebotCommandUsage) String() string {
	var msg string
	for _, s := range ebcu {
		if len(msg) > 0 {
			msg = msg + "\n" + s
		} else {
			msg = s
		}
	}
	return msg
}

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Name() string
	Summary() EvebotCommandSummary
	Usage() EvebotCommandUsage
	OptionalArgs() EvebotCommandArgs
	Examples() EvebotCommandExamples
	Help() *EvebotCommandHelp
	IsHelpRequest(input []string) bool
	IsValidCommand(input []string) bool
	AdditionalArgs(input []string) (EvebotCommandArgs, error)
	ResolveAdditionalArg(argKV []string) EvebotArg
}

func fullHelpMessage() string {
	var msg string

	for _, v := range evebotCommands {
		if len(msg) > 0 {
			msg = msg + "\n" + v.Name() + " help"
		} else {
			msg = v.Name() + " help"
		}

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
