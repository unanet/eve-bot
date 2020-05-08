package api

import (
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

func InitController(c config.Config) []mux.EveController {
	botCommResolver := botcommander.NewResolver()
	eveApiClient := eveapi.NewClient(c.EveAPIConfig)
	slackClient := slack.New(c.SlackUserOauthAccessToken)
	slackProvider := islack.NewProvider(slackClient, botCommResolver, eveApiClient, c.SlackConfig)

	return []mux.EveController{
		NewPingController(),
		NewSlackController(slackProvider),
	}
}
