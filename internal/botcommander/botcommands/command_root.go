package botcommands

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

type EvebotRootCommand struct {
	baseCommand
}

func NewEvebotRootCommand() EvebotRootCommand {
	return EvebotRootCommand{baseCommand{
		name:           "",
		summary:        "Welcome to Eve-Bot! To get started, run:\n```@evebot help```",
		usage:          bothelp.HelpUsage{},
		optionalArgs:   botargs.Args{},
		additionalArgs: botargs.Args{},
		examples:       bothelp.HelpExamples{},
		asyncRequired:  false,
	}}
}

func (cmd EvebotRootCommand) AsyncRequired() bool {
	return cmd.asyncRequired
}

func (cmd EvebotRootCommand) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd EvebotRootCommand) Name() string {
	return cmd.name
}

func (cmd EvebotRootCommand) Help() *bothelp.Help {
	return bothelp.NewEvebotCommandHelp(
		bothelp.EvebotCommandHelpHeaderOpt(cmd.summary.String()),
	)
}

func (cmd EvebotRootCommand) IsHelpRequest() bool {
	return true
}

func (cmd EvebotRootCommand) AdditionalArgs() (botargs.Args, error) {
	return botargs.Args{}, nil
}
