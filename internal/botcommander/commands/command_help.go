package commands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type HelpCmd struct {
	baseCommand
}

func NewHelpCommand(cmdFields []string, channel, user string) EvebotCommand {
	return HelpCmd{baseCommand{
		input:       cmdFields,
		chatDetails: ChatInfo{User: user, Channel: channel},
		name:        "help",
		summary:     "Try running one of the commands below",
		usage: help.Usage{
			"help",
			"{{ command }} help",
		},
		examples:       help.Examples{},
		optionalArgs:   args.Args{},
		requiredParams: params.Params{},
		apiOptions:     make(CommandOptions),
		inputBounds:    InputLengthBounds{Min: 1, Max: -1},
	}}
}

func (cmd HelpCmd) Details() CommandDetails {
	return CommandDetails{
		Name:          cmd.name,
		IsValid:       cmd.ValidInputLength(),
		IsHelpRequest: isHelpRequest(cmd.input, cmd.name),
		AckMsgFn:      baseAckMsg(cmd, cmd.input),
		ErrMsgFn:      cmd.BaseErrMsg(),
	}
}

func (cmd HelpCmd) IsAuthorized(map[string]interface{}, chatChannelInfoFn) bool {
	return true
}

func (cmd HelpCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd HelpCmd) ChatInfo() ChatInfo {
	return cmd.chatDetails
}

func (cmd HelpCmd) Help() *help.Help {
	var nonHelpCmds string
	var nonHelpCmdExamples = help.Examples{}

	for _, v := range nonHelpCmd() {
		if v.Details().Name != cmd.name {
			nonHelpCmds = nonHelpCmds + "\n" + v.Details().Name
			nonHelpCmdExamples = append(nonHelpCmdExamples, v.Details().Name+" help")
		}
	}

	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.CommandsOpt(nonHelpCmds),
		help.UsageOpt(cmd.usage.String()),
		help.ArgsOpt(cmd.optionalArgs.String()),
		help.ExamplesOpt(nonHelpCmdExamples.String()),
	)
}
