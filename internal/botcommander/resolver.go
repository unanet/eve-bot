package botcommander

import (
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
)

// Resolver resolves the input commands and returns a valid EvebotCommand or an Error
type Resolver interface {
	Resolve(input, channel, user string) botcommands.EvebotCommand
}

// EvebotResolver implements the Resolver interface
type EvebotResolver struct{}

// NewResolver instantiates the Resolver
func NewResolver() Resolver {
	return &EvebotResolver{}
}

// Resolve resolves the command input from the Chat User and returns an EvebotCommand
// this is where all of the "magic" happens that basically translates a user command to an EveBot command
func (ebr *EvebotResolver) Resolve(input, channel, user string) botcommands.EvebotCommand {
	// parse the input string and break out into fields (array)
	msgFields := strings.Fields(input)
	if len(msgFields) == 1 {
		// equivalent to just `@evebot`
		// botIDField := msgFields[0]
		return botcommands.NewRootCmd([]string{""}, channel, user)
	}

	cmdFields := msgFields[1:]

	// make sure after you create a new command,
	// you add the New func to the map so that it is picked up here
	newCmdFuncInterface := botcommands.CommandInitializerMap[cmdFields[0]]

	// Make sure the New Command func follows the standard New Command signature
	// =======> func NewCmd(input []string, channel, user, timestamp string) EvebotCommand { }
	if newCmdFuncInterface != nil {
		if newCmdFuncVal, ok := newCmdFuncInterface.(func([]string, string, string) botcommands.EvebotCommand); ok {
			return newCmdFuncVal(cmdFields, channel, user)
		}
		// this is bad - we will want to be alerted on this error
		log.Logger.Error("invalid new command initializer func", zap.String("input", cmdFields[0]))
	}

	return botcommands.NewInvalidCommand(cmdFields, channel, user)
}
