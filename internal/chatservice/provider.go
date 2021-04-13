package chatservice

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/interfaces"

	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/slackservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
)

// ProviderType data structure
type ProviderType string

const (
	// Slack provider type
	Slack ProviderType = "slack"
)

// New returns a chat provider than implements the interface
func New(pt ProviderType, cfg *config.Config) interfaces.ChatProvider {
	switch pt {
	case Slack:
		return slackservice.New(slack.New(cfg.SlackUserOauthAccessToken))
	default:
		return nil
	}
}
