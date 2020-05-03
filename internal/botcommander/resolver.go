package botcommander

import (
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
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

	// make sure after you create new command, you add the New func
	// to the map so that is get's picked up here!
	newCmdFuncInterface := botcommands.CommandInitializerMap[cmdFields[0]]

	// Just make sure your New Command func follows the standard New Command signature
	// =======> func NewCmd(input []string) EvebotCommand { }
	if newCmdFuncInterface != nil {
		if newCmdFuncVal, ok := newCmdFuncInterface.(func([]string) botcommands.EvebotCommand); ok {
			return newCmdFuncVal(cmdFields)
		}
		// this is bad - we will want to be alerted on this error
		log.Logger.Error("invalid new command initializer func", zap.String("input", cmdFields[0]))
	}

	return botcommands.NewInvalidCommand(cmdFields)
}
