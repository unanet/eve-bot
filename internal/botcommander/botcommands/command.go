package botcommands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Name() string
	Help() *bothelp.Help
	IsValid() bool
	IsHelpRequest() bool
	MakeAsyncReq() bool
	AckMsg(userID string) string
	ErrMsg() string
}

//
// @evebot 				(show toplevel evebot help/welcome message)
// @evebot help 		(show toplevel help with full command list and help usage)
// @evebot cmd			(show specific command help)
// @evebot cmd help		(show specific command help)
// @evebot help cmd		(show specific sub/command help)
func isHelpRequest(inputCmd []string, cmdName string) bool {
	if len(inputCmd) == 0 || inputCmd[0] == "help" || inputCmd[len(inputCmd)-1] == "help" || (len(inputCmd) == 1 && inputCmd[0] == cmdName) {
		return true
	}
	return false
}

func baseIsValid(input []string) bool {
	if input == nil || len(input) == 0 {
		return false
	}
	return true
}

func baseAckMsg(cmd EvebotCommand, userID string, cmdInput []string) string {
	if cmd.IsHelpRequest() {
		return fmt.Sprintf("<@%s>...\n\n%s", userID, cmd.Help().String())
	}
	if cmd.IsValid() == false {
		return fmt.Sprintf("Yo <@%s>, one of us goofed up...¯\\_(ツ)_/¯...I don't know what to do with: `%s`\n\nTry running: ```@evebot %s help```\n\n", userID, cmdInput, cmd.Name())
	}
	if len(cmd.ErrMsg()) > 0 {
		return fmt.Sprintf("Whoops <@%s>! I detected some *errors:*\n\n ```%v```", userID, cmd.ErrMsg())
	}
	// Happy Path
	return fmt.Sprintf("Sure <@%s>, I'll `%s` that right away. BRB!", userID, cmd.Name())
}

func baseErrMsg(errs []error) string {
	msg := ""
	if len(errs) > 0 {
		for _, v := range errs {
			if len(msg) == 0 {
				msg = v.Error()
			} else {
				msg = msg + "\n" + v.Error()
			}
		}
	}
	return msg
}

type baseCommand struct {
	input          []string
	name           string
	async          bool
	valid          bool
	errs           []error
	summary        bothelp.Summary
	usage          bothelp.Usage
	optionalArgs   botargs.Args // these are used for the help command
	suppliedArgs   botargs.Args // these are the actual supplied args from the user
	examples       bothelp.Examples
	requiredParams botparams.Params // these are used so we know what the user should supply
	suppliedParams botparams.Params // these are used so we know what the user actually supplied
}
