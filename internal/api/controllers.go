package api

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/executor"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resolver"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/evebotservice"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

// InitController initializes the controller (handlers)
func InitController(cfg *config.Config) []mux.EveController {

	cmdResolver := resolver.NewResolver()
	eveAPI := eveapi.NewClient(cfg.EveAPIConfig)
	chatSvc := chatservice.New(chatservice.Slack, cfg)
	cmdExecutor := executor.NewExecutor(eveAPI, chatSvc)

	svc := evebotservice.New(
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
