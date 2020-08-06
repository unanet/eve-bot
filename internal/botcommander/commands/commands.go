package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

var CommandInitializerMap = map[string]interface{}{
	"help":    NewHelpCommand,
	"deploy":  NewDeployCommand,
	"migrate": NewMigrateCommand,
	"show":    NewShowCommand,
	"set":     NewSetCommand,
	"delete":  NewDeleteCommand,
	"release": NewReleaseCommand,
}

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
	User, Channel string
}

/*
 ------------- Base Command ---------------------
This is the root/abstract/base (i.e. shared/common) struct that all commands use
Ideally

*/
type baseCommand struct {
	input          InputCommand
	inputBounds    InputLengthBounds
	chatDetails    ChatInfo
	name           string
	valid          bool
	errs           []error
	summary        help.Summary
	usage          help.Usage
	examples       help.Examples
	optionalArgs   args.Args
	requiredParams params.Params
	apiOptions     CommandOptions // when we resolve the optionalArgs and requiredParams we hydrate this map for fast lookup
}

func (bc *baseCommand) ValidInputLength() bool {
	return bc.inputBounds.Valid(bc.input)
}

func (bc *baseCommand) ValidMinInputLength() bool {
	return bc.inputBounds.ValidMin(bc.input)
}

func (bc *baseCommand) ValidMaxInputLength() bool {
	return bc.inputBounds.ValidMax(bc.input)
}

func (bc *baseCommand) BaseErrMsg() ErrMsgFn {
	return func() string {
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
}

func baseAckMsg(cmd EvebotCommand, cmdInput []string) ChatAckMsgFn {
	if cmd.Details().IsHelpRequest {
		return ackMsg(fmt.Sprintf("<@%s>...\n\n%s", cmd.ChatInfo().User, cmd.Help().String()), false)
	}
	if cmd.Details().IsValid == false {
		return ackMsg(fmt.Sprintf("Yo <@%s>, one of us goofed up...¯\\_(ツ)_/¯...I don't know what to do with: `%s`\n\nTry running: ```@evebot %s help```\n\n", cmd.ChatInfo().User, cmdInput, cmd.Details().Name), false)
	}
	if len(cmd.Details().ErrMsgFn()) > 0 {
		return ackMsg(fmt.Sprintf("Whoops <@%s>! I detected some command *errors:*\n\n ```%v```", cmd.ChatInfo().User, cmd.Details().ErrMsgFn()), false)
	}
	// Happy Path
	return ackMsg(fmt.Sprintf("Sure <@%s>, I'll `%s` that right away. BRB!", cmd.ChatInfo().User, cmd.Details().Name), true)
}

/*

 ------------- Command Details ---------------------

*/

type ErrMsgFn func() string

type CommandDetails struct {
	Name                   string
	IsValid, IsHelpRequest bool
	ErrMsgFn               ErrMsgFn
	AckMsgFn               ChatAckMsgFn
}

type ChatAckMsgFn func() (string, bool)

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Details() CommandDetails
	Help() *help.Help
	ChatInfo() ChatInfo
	APIOptions() CommandOptions
	IsAuthorized(allowedChannel map[string]interface{}, fn chatChannelInfoFn) bool
}
