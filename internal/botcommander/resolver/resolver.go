package resolver

import (
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"

	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
)

// EvebotResolver implements the Resolver interface
type EvebotResolver struct {
	cmdFactory commands.Factory
}

// New instantiates the Resolver
func New(commandFactory commands.Factory) interfaces.CommandResolver {
	return &EvebotResolver{
		commandFactory,
	}
}

// Resolve resolves the command input from the Chat User and returns an EvebotCommand
// this is where all of the "magic" happens that basically translates a user command to an EveBot command
func (ebr *EvebotResolver) Resolve(input, channel, user string) commands.EvebotCommand {
	// parse the input string and break out into fields (array)
	log.Logger.Info("resolve command", zap.String("input", input))

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
		return commands.NewInvalidCommand(msgFields[1:], channel, user)
	}

	// make sure after you create a new command,
	// you add the New func to the map so that it is picked up here
	if fn := ebr.cmdFactory.Items()[cleanCmdFields[0]]; fn != nil {
		return fn(cleanCmdFields, channel, user)
	}

	log.Logger.Info("invalid command", zap.String("command", cleanCmdFields[0]), zap.String("input", input))
	return commands.NewInvalidCommand(cleanCmdFields, channel, user)
}

func cleanCommandField(cmdFields []string) []string {
	var cleanCmdFields []string
	for _, i := range cmdFields {
		cleanCmdFields = append(cleanCmdFields, commands.CleanUrls(i))
	}
	return cleanCmdFields
}
