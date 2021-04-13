package handlers

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/interfaces"
)

type Factory interface {
	Items() map[string]func(eveAPIClient interfaces.EveAPI, chatSvc interfaces.ChatProvider) CommandHandler
}

type factory struct {
	Map map[string]func(eveAPIClient interfaces.EveAPI, chatSvc interfaces.ChatProvider) CommandHandler
}

func NewFactory() Factory {
	return &factory{
		Map: map[string]func(eveAPIClient interfaces.EveAPI, chatSvc interfaces.ChatProvider) CommandHandler{
			commands.DeployCmdName:  NewDeployHandler,
			commands.ShowCmdName:    NewShowHandler,
			commands.SetCmdName:     NewSetHandler,
			commands.DeleteCmdName:  NewDeleteHandler,
			commands.ReleaseCmdName: NewReleaseHandler,
			commands.RestartCmdName: NewRestartHandler,
			commands.RunCmdName:     NewRunHandler,
		},
	}
}

func (f *factory) Items() map[string]func(eveAPIClient interfaces.EveAPI, chatSvc interfaces.ChatProvider) CommandHandler {
	return f.Map
}
