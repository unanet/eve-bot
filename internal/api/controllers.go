package api

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/evebotservice"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func InitController(cfg config.Config) []mux.EveController {

	cmdResolver := botcommander.NewResolver()
	eveAPI := eveapi.NewClient(cfg.EveAPIConfig)
	chatSvc := chatservice.New(chatservice.Slack, &cfg)
	cmdHandler := botcommander.NewHandler(eveAPI, chatSvc)

	svc := evebotservice.New(
		cfg,
		cmdResolver,
		eveAPI,
		chatSvc,
		cmdHandler,
	)

	return []mux.EveController{
		NewPingController(),
		NewSlackController(svc),
	}
}
