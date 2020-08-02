package commands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

func NewHelpCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultHelpCommand(cmdFields, channel, user)
}

type HelpCmd struct {
	baseCommand
}

func defaultHelpCommand(cmdFields []string, channel, user string) HelpCmd {
	return HelpCmd{baseCommand{
		input:   cmdFields,
		channel: channel,
		user:    user,
		name:    "help",
		summary: "Try running one of the commands below",
		usage: help.Usage{
			"{{ command }} help",
		},
		examples:       help.Examples{},
		optionalArgs:   args.Args{},
		requiredParams: params.Params{},
		apiOptions:     make(CommandOptions),
	}}
}

func (cmd HelpCmd) IsAuthorized(allowedChannel map[string]interface{}, fn chatChannelInfo) bool {
	return true
}

func (cmd HelpCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd HelpCmd) User() string {
	return cmd.user
}

func (cmd HelpCmd) Channel() string {
	return cmd.channel
}

func (cmd HelpCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd HelpCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd HelpCmd) IsValid() bool {
	if len(cmd.errs) > 0 {
		return false
	}
	return baseIsValid(cmd.input)
}

func (cmd HelpCmd) Name() string {
	return cmd.name
}

func (cmd HelpCmd) Help() *help.Help {
	var nonHelpCmds string
	var nonHelpCmdExamples = help.Examples{}

	for _, v := range nonHelpCmd() {
		if v.Name() != cmd.name {
			nonHelpCmds = nonHelpCmds + "\n" + v.Name()
			nonHelpCmdExamples = append(nonHelpCmdExamples, v.Name()+" help")
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

func (cmd HelpCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}
