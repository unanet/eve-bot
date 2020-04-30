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

	// Don't have a clever way to register the commands into this map
	// so make sure after you create new command, you add the New func
	// to the map so that is get's picked up here! Boom!
	newCmdFunc := botcommands.CommandInitializerMap[cmdFields[0]]

	// I know this is magic! And fucking beautiful if you ask me! :)
	// Just make sure your New Command func follows the standard signature
	// =======> func NewCmd(input []string) EvebotCommand { }
	if newCmdFunc != nil {
		return newCmdFunc.(func([]string) botcommands.EvebotCommand)(cmdFields)
	}

	return botcommands.NewInvalidCommand(cmdFields)

}
