package commander

import "fmt"

// Resolver resolves the input commands and returns a valid EvebotCommand or an Error
type Resolver interface {
	Resolve(input []string) (EvebotCommand, error)
}

type EvebotResolver struct{}

func NewResolver() Resolver {
	return &EvebotResolver{}
}

func (ebr *EvebotResolver) Resolve(input []string) (EvebotCommand, error) {

	// This occurs when the user pings evebot without a command
	// example: @evebot
	// thinking about adding an EvebotRootCommand...
	if len(input) <= 0 {
		return NewEvebotHelpCommand(), nil
	}

	for _, v := range evebotCommands {
		if v.Name() == input[0] {
			return v.Initialize(input), nil
		}
	}

	return nil, fmt.Errorf("invalid evebot command: %v", input[0])
}
