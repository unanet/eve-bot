package botcommander

import (
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
)

// Resolver resolves the input commands and returns a valid EvebotCommand or an Error
type Resolver interface {
	Resolve(input string) botcommands.EvebotCommand
}

// EvebotResolver implements the Resolver interface
type EvebotResolver struct{}

// NewResolver instantiates the Resolver
func NewResolver() Resolver {
	return &EvebotResolver{}
}

// Resolve resolves the input command and returns a valid EvebotCommand or an error
func (ebr *EvebotResolver) Resolve(input string) botcommands.EvebotCommand {
	// parse the input string
	msgFields := strings.Fields(input)
	// equivalent to just `@evebot`
	if len(msgFields) == 1 {
		// botIDField := msgFields[0]
		return botcommands.NewRootCmd()
	}

	cmdFields := msgFields[1:]
	cmdName := cmdFields[0]

	switch strings.ToLower(cmdName) {
	case "deploy":
		return botcommands.NewDeployCommand(cmdFields)
	case "help":
		return botcommands.NewHelpCommand(cmdFields)
	case "migrate":
		return botcommands.NewMigrateCommand(cmdFields)
	default:
		return botcommands.NewInvalidCommand(cmdFields)
	}

}
