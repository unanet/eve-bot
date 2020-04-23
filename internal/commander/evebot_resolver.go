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
	if len(input) <= 1 {
		return nil, fmt.Errorf("invalid evebot command")
	}

	for _, v := range evebotCommands {
		if v.Name() == input[0] {
			return v, nil
		}
	}

	return nil, fmt.Errorf("invalid evebot command: %v", input[0])
}
