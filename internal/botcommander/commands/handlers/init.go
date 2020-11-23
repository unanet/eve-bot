package handlers

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
)

var (
	// CommandHandlerMap maps command to handlers
	CommandHandlerMap = map[string]interface{}{
		commands.DeployCmdName:  NewDeployHandler,
		commands.MigrateCmdName: NewMigrateHandler,
		commands.ShowCmdName:    NewShowHandler,
		commands.SetCmdName:     NewSetHandler,
		commands.DeleteCmdName:  NewDeleteHandler,
		commands.ReleaseCmdName: NewReleaseHandler,
		commands.RestartCmdName: NewRestartHandler,
		commands.RunCmdName:     NewRunHandler,
	}
)
