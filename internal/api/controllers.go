package api

import (
	"github.com/go-chi/chi"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands/handlers"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/executor"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resolver"
	chat "gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/service"
)

type Controller interface {
	Setup(chi.Router)
}

// InitController initializes the controller (handlers)
func InitController(cfg *config.Config) []Controller {

	cmdResolver := resolver.New(commands.NewFactory())

	eveAPI := eveapi.New(cfg.EveAPIConfig)
	chatSvc := chat.New(chat.Slack, cfg)
	cmdExecutor := executor.New(eveAPI, chatSvc, handlers.NewFactory())

	svc := service.New(
		cfg,
		cmdResolver,
		eveAPI,
		chatSvc,
		cmdExecutor,
	)

	return []Controller{
		NewPingController(),
		NewSlackController(svc),
		NewEveController(svc),
	}
}
