package chatservice

import (
	"context"

	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/slackservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
)

type ProviderType string

const (
	Slack ProviderType = "slack"
)

type Provider interface {
	PostMessage(ctx context.Context, msg, channel string)
	PostMessageThread(ctx context.Context, msg, channel, ts string) (timestamp string)
	ErrorNotification(ctx context.Context, user, channel string, err error)
	ErrorNotificationThread(ctx context.Context, user, channel, ts string, err error)
	UserNotificationThread(ctx context.Context, msg, user, channel, ts string)
	DeploymentNotificationThread(ctx context.Context, msg, user, channel, ts string)
	GetUser(ctx context.Context, user string) (*chatmodels.ChatUser, error)
	PostLinkMessageThread(ctx context.Context, msg string, user string, channel string, ts string)
}

func New(st ProviderType, cfg *config.Config) Provider {
	switch st {
	case Slack:
		return slackservice.New(slack.New(cfg.SlackUserOauthAccessToken))
	default:
		return nil
	}
}
