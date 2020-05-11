package slack

import (
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

var (
	// Setting this up here for now
	// no need to resolve this every request
	callBackURL   string
	signingSecret string
)

// Provider provides access to the Slack Client
// and the deps required for this package
type Provider struct {
	Client          *slack.Client
	CommandResolver botcommander.Resolver
	EveAPIClient    eveapi.Client
	Cfg             Config
}

// NewProvider creates a new provider
func NewProvider(sClient *slack.Client, commander botcommander.Resolver, eveAPIClient eveapi.Client, cfg Config) *Provider {
	callBackURL = eveAPIClient.CallBackURL()
	signingSecret = cfg.SlackSigningSecret
	return &Provider{
		Client:          sClient,
		CommandResolver: commander,
		EveAPIClient:    eveAPIClient,
		Cfg:             cfg,
	}
}
