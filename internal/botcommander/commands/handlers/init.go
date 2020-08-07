package handlers

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
)

var (
	// Deploy and Migrate share the same command
	// the api honors the same request/response signature for both
	CommandHandlerMap = map[string]interface{}{
		commands.DeployCmdName:  NewDeployHandler,
		commands.MigrateCmdName: NewMigrateHandler,
		commands.ShowCmdName:    NewShowHandler,
		commands.SetCmdName:     NewSetHandler,
		commands.DeleteCmdName:  NewDeleteHandler,
		commands.ReleaseCmdName: NewReleaseHandler,
	}
)
