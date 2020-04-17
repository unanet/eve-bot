package api

import (
	authController "gitlab.unanet.io/devops/eve-bot/internal/api/controller/auth"
	pingController "gitlab.unanet.io/devops/eve-bot/internal/api/controller/ping"
	slackController "gitlab.unanet.io/devops/eve-bot/internal/api/controller/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

// Controllers contains all of the route controller (handlers)
var Controllers = []mux.EveController{
	pingController.New(),
	slackController.New(islack.NewProvider(config.Values().SlackConfig)),
	authController.New(),
}
