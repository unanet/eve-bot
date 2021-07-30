package api

import (
	"github.com/go-chi/chi"
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/commands/handlers"
	"github.com/unanet/eve-bot/internal/botcommander/executor"
	"github.com/unanet/eve-bot/internal/botcommander/resolver"
	chat "github.com/unanet/eve-bot/internal/chatservice"
	"github.com/unanet/eve-bot/internal/config"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/eve-bot/internal/service"
)

type Controller interface {
	Setup(chi.Router)
}

// initController initializes the controller (handlers)
func initController(cfg *config.Config) []Controller {

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
