package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"
)

type chatChannelInfoFn func(string) (chatmodels.Channel, error)

type CommandOptions map[string]interface{}

/*

 ------------- Command Input Validation ---------------------

*/
type InputLengthBounds struct {
	Min, Max int
}

func (ilb *InputLengthBounds) ValidMax(input []string) bool {
	if ilb.Max <= 0 {
		return true
	}
	return len(input) <= ilb.Max
}

func (ilb *InputLengthBounds) ValidMin(input []string) bool {
	return len(input) >= ilb.Min
}

func (ilb *InputLengthBounds) Valid(input []string) bool {
	return ilb.ValidMax(input) && ilb.ValidMin(input)
}

type InputCommand []string

func (ic InputCommand) Length() int {
	return len(ic)
}

type ChatInfo struct {
	User, Channel, CommandName string
}

/*
 ------------- Base Command ---------------------
This is the root/abstract/base (i.e. shared/common) struct that all commands share
*/
type baseCommand struct {
	input      InputCommand
	bounds     InputLengthBounds
	info       ChatInfo
	valid      bool
	errs       []error
	arguments  args.Args
	parameters params.Params
	opts       CommandOptions // when we resolve the arguments and parameters we hydrate this map for fast lookup
}

func (bc *baseCommand) IsHelpRequest() bool {
	if len(bc.input) == 0 || bc.input[0] == "help" || bc.input[len(bc.input)-1] == "help" || (len(bc.input) == 1 && bc.input[0] == bc.info.CommandName) {
		return true
	}
	return false
}

func (bc *baseCommand) ValidInputLength() bool {
	return bc.bounds.Valid(bc.input)
}

func (bc *baseCommand) ValidMinInputLength() bool {
	return bc.bounds.ValidMin(bc.input)
}

func (bc *baseCommand) ValidMaxInputLength() bool {
	return bc.bounds.ValidMax(bc.input)
}

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

func (bc *baseCommand) BaseAckMsg(cmdHelp string) (string, bool) {
	if bc.IsHelpRequest() {
		return fmt.Sprintf("<@%s>...\n\n%s", bc.info.User, cmdHelp), false
	}
	if bc.ValidInputLength() == false {
		return fmt.Sprintf("Yo <@%s>, one of us goofed up...¯\\_(ツ)_/¯...I don't know what to do with: `%s`\n\nTry running: ```@evebot %s help```\n\n", bc.info.User, bc.input, bc.info.CommandName), false
	}
	if len(bc.BaseErrMsg()) > 0 {
		return fmt.Sprintf("Whoops <@%s>! I detected some command *errors:*\n\n ```%v```", bc.info.User, bc.BaseErrMsg()), false
	}
	// Happy Path
	return fmt.Sprintf("Sure <@%s>, I'll `%s` that right away. BRB!", bc.info.User, bc.info.CommandName), true
}

// EvebotCommand interface (each evebot command needs to implement this interface)
type EvebotCommand interface {
	ChatInfo() ChatInfo
	DynamicOptions() CommandOptions
	AckMsg() (string, bool)
	IsAuthorized(allowedChannel map[string]interface{}, fn chatChannelInfoFn) bool
}
