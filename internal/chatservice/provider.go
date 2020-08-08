package chatservice

import (
	"context"

	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/slackservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
)

// ProviderType data structure
type ProviderType string

const (
	// Slack provider type
	Slack ProviderType = "slack"
)

// Provider represents a Chat Provider Interface
type Provider interface {
	GetChannelInfo(ctx context.Context, channelID string) (chatmodels.Channel, error)
	PostMessage(ctx context.Context, msg, channel string) (timestamp string)
	PostMessageThread(ctx context.Context, msg, channel, ts string) (timestamp string)
	ErrorNotification(ctx context.Context, user, channel string, err error)
	ErrorNotificationThread(ctx context.Context, user, channel, ts string, err error)
	UserNotificationThread(ctx context.Context, msg, user, channel, ts string)
	DeploymentNotificationThread(ctx context.Context, msg, user, channel, ts string)
	GetUser(ctx context.Context, user string) (*chatmodels.ChatUser, error)
	PostLinkMessageThread(ctx context.Context, msg string, user string, channel string, ts string)
	ShowResultsMessageThread(ctx context.Context, msg, user, channel, ts string)
	ReleaseResultsMessageThread(ctx context.Context, msg, user, channel, ts string)
}

// New returns a chat provider than implements the interface
func New(st ProviderType, cfg *config.Config) Provider {
	switch st {
	case Slack:
		return slackservice.New(slack.New(cfg.SlackUserOauthAccessToken))
	default:
		return nil
	}
}
