package commander

func init() {
	// Add all the Evebot commands here on init
	evebotCommands = []EvebotCommand{
		NewEvebotDeployCommand(),
	}

	// Set the Full Help Message Once on init
	CmdHelpMsgs = fullHelpMessage()
}

var (
	evebotCommands []EvebotCommand
	CmdHelpMsgs    string
)

type EvebotCommandExamples []string

func (ebce EvebotCommandExamples) HelpMsg() string {
	var helpMsg string
	for _, v := range ebce {
		if len(helpMsg) > 0 {
			helpMsg = helpMsg + "\n" + v
		} else {
			helpMsg = v
		}
	}
	return helpMsg
}

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Name() string
	Examples() EvebotCommandExamples
	IsHelpRequest(input []string) bool
	IsValidCommand(input []string) bool
	AdditionalArgs(input []string) (EvebotArgs, error)
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
