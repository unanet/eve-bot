package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

type RootCmd struct {
	baseCommand
}

func NewRootCmd() RootCmd {
	return RootCmd{baseCommand{
		name:           "",
		summary:        "Welcome to `@evebot`! To get started, run:\n```@evebot help```",
		usage:          bothelp.Usage{},
		optionalArgs:   botargs.Args{},
		additionalArgs: botargs.Args{},
		examples:       bothelp.Examples{},
		asyncRequired:  false,
	}}
}

func (cmd RootCmd) IsValid() bool {
	if cmd.input == nil {
		return false
	}
	return true
}

func (cmd RootCmd) AsyncRequired() bool {
	return cmd.asyncRequired
}

func (cmd RootCmd) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd RootCmd) Name() string {
	return cmd.name
}

func (cmd RootCmd) Help() *bothelp.Help {
	return bothelp.New(
		bothelp.HeaderOpt(cmd.summary.String()),
	)
}

func (cmd RootCmd) IsHelpRequest() bool {
	return true
}

func (cmd RootCmd) AdditionalArgs() (botargs.Args, error) {
	return botargs.Args{}, nil
}
