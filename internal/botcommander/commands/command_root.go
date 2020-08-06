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
		chatDetails:    ChatInfo{User: user, Channel: channel},
		name:           "",
		summary:        "Welcome to `@evebot`! To get started, run:\n```@evebot help```",
		usage:          help.Usage{},
		examples:       help.Examples{},
		optionalArgs:   args.Args{},
		requiredParams: params.Params{},
		apiOptions:     make(CommandOptions),
	}}
}

func (cmd RootCmd) Details() CommandDetails {
	return CommandDetails{
		Name:          cmd.name,
		IsValid:       true,
		IsHelpRequest: isHelpRequest(cmd.input, cmd.name),
		AckMsgFn:      baseAckMsg(cmd, cmd.input),
		ErrMsgFn:      cmd.BaseErrMsg(),
	}
}

func (cmd RootCmd) IsAuthorized(allowedChannel map[string]interface{}, fn chatChannelInfoFn) bool {
	return true
}

func (cmd RootCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd RootCmd) ChatInfo() ChatInfo {
	return cmd.chatDetails
}

func (cmd RootCmd) Help() *help.Help {
	return help.New(help.HeaderOpt(cmd.summary.String()))
}
