package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type InvalidCmd struct {
	baseCommand
}

func NewInvalidCommand(cmdFields []string, channel, user string) EvebotCommand {
	return InvalidCmd{baseCommand{
		input:          cmdFields,
		chatDetails:    ChatInfo{User: user, Channel: channel},
		name:           "",
		summary:        help.Summary(fmt.Sprintf("I don't know how to execute the `%s` command.\n\nTry running: ```@evebot help```\n", cmdFields)),
		usage:          help.Usage{},
		examples:       help.Examples{},
		optionalArgs:   args.Args{},
		requiredParams: params.Params{},
		apiOptions:     make(CommandOptions),
		inputBounds:    InputLengthBounds{Min: 1, Max: -1},
	}}
}

func (cmd InvalidCmd) Details() CommandDetails {
	return CommandDetails{
		Name:          cmd.name,
		IsValid:       cmd.ValidInputLength(),
		IsHelpRequest: isHelpRequest(cmd.input, cmd.name),
		AckMsgFn:      baseAckMsg(cmd, cmd.input),
		ErrMsgFn:      cmd.BaseErrMsg(),
	}
}

func (cmd InvalidCmd) IsAuthorized(map[string]interface{}, chatChannelInfoFn) bool {
	return true
}

func (cmd InvalidCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd InvalidCmd) ChatInfo() ChatInfo {
	return cmd.chatDetails
}

func (cmd InvalidCmd) Help() *help.Help {

	var nonHelpCmds string
	var nonHelpCmdExamples = help.Examples{}

	for _, v := range nonHelpCmd() {
		if v.Details().Name != "help" {
			nonHelpCmds = nonHelpCmds + "\n" + v.Details().Name
			nonHelpCmdExamples = append(nonHelpCmdExamples, v.Details().Name+" help")
		}
	}

	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.CommandsOpt(nonHelpCmds),
		help.ExamplesOpt(nonHelpCmdExamples.String()),
	)

}
