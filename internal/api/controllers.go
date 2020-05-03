package api

import (
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

// Controllers contains all of the route controller (handlers)
var Controllers = []mux.EveController{
	NewPingController(),
	NewSlackController(
		islack.NewProvider(
			slack.New(config.Values().SlackUserOauthAccessToken),
			botcommander.NewResolver(),
			eveapi.NewClient(config.Values().EveAPIConfig),
			config.Values().SlackConfig,
		),
	),
}
