package resolver

import (
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
)

// Resolver resolves the input commands and returns a valid EvebotCommand or an Error
type Resolver interface {
	Resolve(input, channel, user string) commands.EvebotCommand
}

// EvebotResolver implements the Resolver interface
type EvebotResolver struct{}

// NewResolver instantiates the Resolver
func NewResolver() Resolver {
	return &EvebotResolver{}
}

// Resolve resolves the command input from the Chat User and returns an EvebotCommand
// this is where all of the "magic" happens that basically translates a user command to an EveBot command
func (ebr *EvebotResolver) Resolve(input, channel, user string) commands.EvebotCommand {
	// parse the input string and break out into fields (array)
	msgFields := strings.Fields(input)
	if len(msgFields) == 1 {
		// equivalent to just `@evebot`
		// botIDField := msgFields[0]
		return commands.NewRootCmd([]string{""}, channel, user)
	}

	// scrub the input fields for invalid data (link encoding)
	cleanCmdFields := cleanCommandField(msgFields[1:])
	if cleanCmdFields == nil {
		log.Logger.Error("invalid clean cmd fields")
		return commands.NewInvalidCommand(cleanCmdFields, channel, user)
	}

	// make sure after you create a new command,
	// you add the New func to the map so that it is picked up here
	newCmdFuncInterface := commands.CommandInitializerMap[cleanCmdFields[0]]
	if newCmdFuncInterface == nil {
		// this is bad - we will want to be alerted on this error
		log.Logger.Error("invalid new command initializer func", zap.String("input", cleanCmdFields[0]))
		return commands.NewInvalidCommand(cleanCmdFields, channel, user)
	}

	// Make sure the New Command func follows the standard New Command signature
	// =======> func NewCmd(input []string, channel, user string) EvebotCommand { }
	if newCmdFuncVal, ok := newCmdFuncInterface.(func([]string, string, string) commands.EvebotCommand); ok {
		return newCmdFuncVal(cleanCmdFields, channel, user)
	}

	// this is bad - we will want to be alerted on this error
	log.Logger.Error("unknown cmd resolved", zap.String("input", cleanCmdFields[0]))
	return commands.NewInvalidCommand(cleanCmdFields, channel, user)
}

func cleanCommandField(cmdFields []string) []string {
	var cleanCmdFields []string
	for _, i := range cmdFields {
		cleanCmdFields = append(cleanCmdFields, commands.CleanUrls(i))
	}
	return cleanCmdFields
}
