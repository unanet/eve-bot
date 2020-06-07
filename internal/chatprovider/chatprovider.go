package chatprovider

import (
	"context"

	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/chatprovider/slackchatprovider"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
)

type ServiceType string

const (
	Slack ServiceType = "slack"
)

type Service interface {
	PostMessage(ctx context.Context)
}

func New(st ServiceType, cfg *config.Config) Service {
	switch st {
	case Slack:
		return slackchatprovider.New(slack.New(cfg.SlackUserOauthAccessToken))
	default:
		return nil
	}
}
