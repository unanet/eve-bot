package commands

import (
	"github.com/unanet/eve-bot/internal/botcommander/args"
	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
)

type rootCmd struct {
	baseCommand
}

const (
	rootCmdHelpSummary = help.Summary("Welcome to `@evebot`! To get started, run:\n```@evebot help```")
)

// NewRootCmd creates a New RootCmd that implements the EvebotCommand interface
func NewRootCmd(cmdFields []string, channel, user string) EvebotCommand {
	cmd := rootCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: ""},
		arguments:  args.Args{},
		parameters: params.Params{},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 0, Max: -1},
	}}
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd rootCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(help.HeaderOpt(rootCmdHelpSummary.String())).String())
}

// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd rootCmd) IsAuthorized(map[string]interface{}, chatChannelInfoFn, chatUserInfoFn) bool {
	return true
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd rootCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd rootCmd) Info() ChatInfo {
	return cmd.info
}
