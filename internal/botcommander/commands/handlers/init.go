package handlers

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

type Factory interface {
	Items() map[string]func(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler
}

type factory struct {
	Map map[string]func(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler
}

func NewFactory() Factory {
	return &factory{
		Map: map[string]func(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler{
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

func (f *factory) Items() map[string]func(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return f.Map
}
