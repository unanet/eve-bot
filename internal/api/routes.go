package api

import (
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

// Controllers contains all of the route controller (handlers)
var Controllers = []mux.EveController{
	NewPingController(),
	NewSlackController(islack.NewProvider(config.Values().SlackConfig)),
}
