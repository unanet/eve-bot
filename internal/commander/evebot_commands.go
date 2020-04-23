package commander

import (
	"fmt"
)

// Add all the Evebot commands here
// As we add more commands this will need to be updated
func init() {
	evebotCommands = []EvebotCommand{
		NewEvebotDeployCommand(),
	}

	// Set the Full Help Message Once
	CmdHelpMsgs = fullHelpMessage()
}

var (
	evebotCommands []EvebotCommand
	CmdHelpMsgs    string
)

type EvebotCommandExamples []string

func (ebce EvebotCommandExamples) HelpMsg() string {
	var helpMsg string
	for _, v := range ebce {
		if len(helpMsg) > 0 {
			helpMsg = helpMsg + "\n" + v
		} else {
			helpMsg = v
		}
	}
	return helpMsg
}

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Name() string
	Examples() EvebotCommandExamples
	IsHelpRequest(input []string) bool
}

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

func fullHelpMessage() string {
	var msg string

	for _, v := range evebotCommands {
		if len(msg) > 0 {
			msg = msg + "\n" + v.Name() + " help"
		} else {
			msg = v.Name() + " help"
		}

	}

	return msg
}
