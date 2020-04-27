package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Name() string
	Help() *bothelp.Help
	Initialize(input []string) EvebotCommand
	IsHelpRequest() bool
	AdditionalArgs() (botargs.Args, error)
	AsyncRequired() bool
}

func isHelpRequest(inputCmd []string, cmdName string) bool {
	if len(inputCmd) == 0 || inputCmd[0] == "help" || inputCmd[len(inputCmd)-1] == "help" || (len(inputCmd) == 1 && inputCmd[0] == cmdName) {
		return true
	}
	return false
}

type baseCommand struct {
	input          []string
	name           string
	asyncRequired  bool
	summary        bothelp.HelpSummary
	usage          bothelp.HelpUsage
	optionalArgs   botargs.Args // these are used for the help command
	additionalArgs botargs.Args // these are the actual supplied args from the user
	examples       bothelp.HelpExamples
}
