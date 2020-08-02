package commands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type RootCmd struct {
	baseCommand
}

func NewRootCmd(cmdFields []string, channel, user string) EvebotCommand {
	return RootCmd{baseCommand{
		input:          cmdFields,
		channel:        channel,
		user:           user,
		name:           "",
		summary:        "Welcome to `@evebot`! To get started, run:\n```@evebot help```",
		usage:          help.Usage{},
		examples:       help.Examples{},
		optionalArgs:   args.Args{},
		requiredParams: params.Params{},
		apiOptions:     make(CommandOptions),
	}}
}

func (cmd RootCmd) IsAuthorized(allowedChannel map[string]interface{}, fn chatChannelInfo) bool {
	return true
}

func (cmd RootCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd RootCmd) User() string {
	return cmd.user
}

func (cmd RootCmd) Channel() string {
	return cmd.channel
}

func (cmd RootCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd RootCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd RootCmd) IsValid() bool {
	return true
}

func (cmd RootCmd) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd RootCmd) Name() string {
	return cmd.name
}

func (cmd RootCmd) Help() *help.Help {
	return help.New(help.HeaderOpt(cmd.summary.String()))
}

func (cmd RootCmd) IsHelpRequest() bool {
	return true
}
