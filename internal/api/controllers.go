package api

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands/handlers"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/executor"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resolver"
	chat "gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

// InitController initializes the controller (handlers)
func InitController(cfg *config.Config) []mux.EveController {

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

	return []mux.EveController{
		NewPingController(),
		NewSlackController(svc),
		NewEveController(svc),
	}
}
