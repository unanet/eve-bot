package botcommander

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
)

// Resolver resolves the input commands and returns a valid EvebotCommand or an Error
type Resolver interface {
	Resolve(input []string) (botcommands.EvebotCommand, error)
}

// EvebotResolver implements the Resolver interface
type EvebotResolver struct{}

// NewResolver instantiates the Resolver
func NewResolver() Resolver {
	return &EvebotResolver{}
}

// Resolve resolves the input command and returns a valid EvebotCommand or an error
func (ebr *EvebotResolver) Resolve(input []string) (botcommands.EvebotCommand, error) {
	// This occurs when the user pings evebot without a command
	// example: @evebot
	// thinking about adding an RootCmd...
	if len(input) <= 0 {
		return botcommands.NewRootCmd(), nil
	}

	// Match the command input with a command name
	for _, v := range botcommands.EvebotCommands {
		if v.Name() == input[0] {
			return v.Initialize(input), nil
		}
	}

	// Didn't find a match for the command
	return nil, fmt.Errorf("invalid evebot command: %v", input[0])
}
