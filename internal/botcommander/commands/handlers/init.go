package handlers

import (
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/service"
)

type Factory interface {
	Items() map[string]func(svc *service.Provider) CommandHandler
}

type factory struct {
	Map map[string]func(svc *service.Provider) CommandHandler
}

func NewFactory() Factory {
	return &factory{
		Map: map[string]func(svc *service.Provider) CommandHandler{
			commands.DeployCmdName:           NewDeployHandler,
			commands.ShowCmdName:             NewShowHandler,
			commands.SetCmdName:              NewSetHandler,
			commands.DeleteCmdName:           NewDeleteHandler,
			commands.ReleaseArtifactCmdName:  NewReleaseArtifactHandler,
			commands.ReleaseNamespaceCmdName: NewReleaseNamespaceHandler,
			commands.RestartCmdName:          NewRestartHandler,
			commands.RunCmdName:              NewRunHandler,
			commands.AuthCmdName:             NewAuthHandler,
		},
	}
}

func (f *factory) Items() map[string]func(svc *service.Provider) CommandHandler {
	return f.Map
}
