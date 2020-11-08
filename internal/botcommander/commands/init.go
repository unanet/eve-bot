package commands

import (
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
)

var (
	// CommandInitializerMap is the main map that holds all commands
	CommandInitializerMap = map[string]interface{}{
		helpCmdName:    NewHelpCommand,
		DeployCmdName:  NewDeployCommand,
		MigrateCmdName: NewMigrateCommand,
		ShowCmdName:    NewShowCommand,
		SetCmdName:     NewSetCommand,
		DeleteCmdName:  NewDeleteCommand,
		ReleaseCmdName: NewReleaseCommand,
		RestartCmdName: NewRestartCommand,
	}
	// NonHelpCommandExamples is hydrated during init and holds all of the non-helper command examples
	NonHelpCommandExamples = help.Examples{}
	// NonHelpCmds holds all of the non-helper command names
	NonHelpCmds string
)

func init() {
	// Iterate the full command map and extract the Non-Help Command
	// we utilize these for system wide help calls
	for k, v := range CommandInitializerMap {
		if k != helpCmdName {
			nonHelpCmd := v.(func([]string, string, string) EvebotCommand)([]string{}, "", "")
			NonHelpCmds = NonHelpCmds + "\n" + nonHelpCmd.Info().CommandName
			NonHelpCommandExamples = append(NonHelpCommandExamples, fmt.Sprintf("%s %s", nonHelpCmd.Info().CommandName, helpCmdName))
		}
	}
}
