package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

func NewInvalidCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultInvalidCommand(cmdFields, channel, user)
}

type InvalidCmd struct {
	baseCommand
}

func defaultInvalidCommand(cmdFields []string, channel, user string) InvalidCmd {
	return InvalidCmd{baseCommand{
		input:          cmdFields,
		channel:        channel,
		user:           user,
		name:           "",
		summary:        help.Summary(fmt.Sprintf("I don't know how to execute the `%s` command.\n\nTry running: ```@evebot help```\n", cmdFields)),
		usage:          help.Usage{},
		examples:       help.Examples{},
		optionalArgs:   args.Args{},
		requiredParams: params.Params{},
		apiOptions:     make(CommandOptions),
	}}
}

func (cmd InvalidCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd InvalidCmd) User() string {
	return cmd.user
}

func (cmd InvalidCmd) Channel() string {
	return cmd.channel
}

func (cmd InvalidCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd InvalidCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd InvalidCmd) IsValid() bool {
	return false
}

func (cmd InvalidCmd) Name() string {
	return cmd.name
}

func (cmd InvalidCmd) Help() *help.Help {

	var nonHelpCmds string
	var nonHelpCmdExamples = help.Examples{}

	for _, v := range nonHelpCmd() {
		if v.Name() != "help" {
			nonHelpCmds = nonHelpCmds + "\n" + v.Name()
			nonHelpCmdExamples = append(nonHelpCmdExamples, v.Name()+" help")
		}
	}

	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.CommandsOpt(nonHelpCmds),
		help.ExamplesOpt(nonHelpCmdExamples.String()),
	)

}

func (cmd InvalidCmd) IsHelpRequest() bool {
	return true
}
