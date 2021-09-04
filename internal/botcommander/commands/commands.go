package commands

import (
	"context"
	"fmt"

	"github.com/unanet/eve-bot/internal/botcommander/args"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/botcommander/resources"
	"github.com/unanet/eve-bot/internal/chatservice/chatmodels"
)

// CommandOptions is the dynamic command options that become Key:Value pairs
type CommandOptions map[string]interface{}

// InputLengthBounds structures the Min/Max
type InputLengthBounds struct {
	Min, Max int
}

// ValidMax verifies that the input is less than or equal to the defined Max value
func (ilb *InputLengthBounds) ValidMax(input []string) bool {
	if ilb.Max <= 0 {
		return true
	}
	return len(input) <= ilb.Max
}

// ValidMin verifies that the input is greater than or equal to the defined Min value
func (ilb *InputLengthBounds) ValidMin(input []string) bool {
	return len(input) >= ilb.Min
}

// Valid verifies that the input meets the minmum/maximum lengths
func (ilb *InputLengthBounds) Valid(input []string) bool {
	return ilb.ValidMax(input) && ilb.ValidMin(input)
}

// ChatInfo contains the Chat Command Info
type ChatInfo struct {
	User, Channel, CommandName string
}

// baseCommand
// the root/abstract/base (i.e. shared/common) struct that all commands share
type baseCommand struct {
	input      []string
	bounds     InputLengthBounds
	info       ChatInfo
	errs       []error
	arguments  args.Args
	parameters params.Params
	opts       CommandOptions // when we resolve the arguments and parameters we hydrate this map for fast lookup
}

func (bc *baseCommand) verifyInput() {
	if !bc.ValidInputLength() {
		bc.errs = append(bc.errs, fmt.Errorf("invalid input length: %v", bc.input))
	}
}

func (bc *baseCommand) initializeResource() {
	if len(bc.input) < 2 {
		bc.errs = append(bc.errs, fmt.Errorf("invalid input: %v", bc.input))
		return
	}
	if ok := resources.FullResourceMap[bc.input[1]]; ok {
		bc.opts["resource"] = bc.input[1]
	} else {
		bc.errs = append(bc.errs, fmt.Errorf("invalid requested resource: %v", bc.input))
		return
	}

	if bc.opts["resource"] == nil {
		bc.errs = append(bc.errs, fmt.Errorf("invalid resource: %v", bc.input))
		return
	}
}

// IsHelpRequest checks if the command is a request for help
func (bc *baseCommand) IsHelpRequest() bool {
	// There is no help for auth
	// @evebot auth
	if bc.info.CommandName == AuthCmdName {
		return false
	}
	if len(bc.input) == 0 ||
		bc.input[0] == helpCmdName ||
		bc.input[len(bc.input)-1] == helpCmdName ||
		(len(bc.input) == 1 && bc.input[0] == bc.info.CommandName) {
		return true
	}
	return false
}

// ValidInputLength checks the input length
func (bc *baseCommand) ValidInputLength() bool {
	return bc.bounds.Valid(bc.input)
}

// ValidMinInputLength validates the minimum length
func (bc *baseCommand) ValidMinInputLength() bool {
	return bc.bounds.ValidMin(bc.input)
}

// ValidMaxInputLength validates the maximum length
func (bc *baseCommand) ValidMaxInputLength() bool {
	return bc.bounds.ValidMax(bc.input)
}

// BaseErrMsg converts the err slice to a string
func (bc *baseCommand) BaseErrMsg() string {
	msg := ""
	if len(bc.errs) > 0 {
		for _, v := range bc.errs {
			if len(msg) == 0 {
				msg = v.Error()
			} else {
				msg = msg + "\n" + v.Error()
			}
		}
	}
	return msg
}

// BaseAckMsg is figures out what to send to the user immediately (acknowledgment)
func (bc *baseCommand) BaseAckMsg(cmdHelp string) (string, bool) {
	if bc.IsHelpRequest() || bc.info.CommandName == "" {
		return fmt.Sprintf("<@%s>...\n\n%s", bc.info.User, cmdHelp), false
	}
	if !bc.ValidInputLength() {
		return fmt.Sprintf("Yo <@%s>, one of us goofed up...¯\\_(ツ)_/¯...I don't know what to do with: `%s`\n\nTry running: ```@evebot %s help```\n\n", bc.info.User, bc.input, bc.info.CommandName), false
	}
	if len(bc.BaseErrMsg()) > 0 {
		return fmt.Sprintf("Whoops <@%s>! I detected some command *errors:*\n\n ```%v```", bc.info.User, bc.BaseErrMsg()), false
	}
	// Happy Path
	return fmt.Sprintf("Sure <@%s>, I'll `%s` that right away. BRB!", bc.info.User, bc.info.CommandName), true
}

type ChatChannelInfoFn func(context.Context, string) (chatmodels.Channel, error)

// EvebotCommand interface (each evebot command needs to implement this interface)
type EvebotCommand interface {
	Info() ChatInfo
	Options() CommandOptions
	AckMsg() (string, bool)
}
