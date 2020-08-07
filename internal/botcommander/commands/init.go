package commands

import "gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"

var (
	CommandInitializerMap  = map[string]interface{}{}
	NonHelpCommands        []EvebotCommand
	NonHelpCommandExamples help.Examples
	NonHelpCmds            string
)

func init() {
	CommandInitializerMap = map[string]interface{}{
		helpCmdName:    NewHelpCommand,
		DeployCmdName:  NewDeployCommand,
		MigrateCmdName: NewMigrateCommand,
		ShowCmdName:    NewShowCommand,
		SetCmdName:     NewSetCommand,
		DeleteCmdName:  NewDeleteCommand,
		ReleaseCmdName: NewReleaseCommand,
	}

	NonHelpCommands = []EvebotCommand{}

	for k, v := range CommandInitializerMap {
		if k != "help" {
			NonHelpCommands = append(NonHelpCommands, v.(func([]string, string, string) EvebotCommand)([]string{}, "", ""))
		}
	}

	NonHelpCommandExamples = help.Examples{}

	for _, v := range NonHelpCommands {
		if v.ChatInfo().CommandName != helpCmdName {
			NonHelpCmds = NonHelpCmds + "\n" + v.ChatInfo().CommandName
			NonHelpCommandExamples = append(NonHelpCommandExamples, v.ChatInfo().CommandName+" help")
		}
	}

}
