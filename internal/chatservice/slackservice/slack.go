package slackservice

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
	"github.com/unanet/eve-bot/internal/chatservice/chatmodels"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

// Provider is the Slack provider which wraps the slack the client
type Provider struct {
	client *slack.Client
}

// New returns a new Slack provider
func New(c *slack.Client) Provider {
	return Provider{client: c}
}

func (sp Provider) handleDevOpsErrorNotification(ctx context.Context, err error) {
	if err != nil {
		log.Logger.Error("critical devops error", zap.Error(err))
		_, _, _ = sp.client.PostMessageContext(ctx, devOpsMonitoringChannel, slack.MsgOptionText(errMessage(err), false))
	}
}

// GetChannelInfo returns the slack channel info
func (sp Provider) GetChannelInfo(ctx context.Context, channelID string) (chatmodels.Channel, error) {
	slackChannel, err := sp.client.GetConversationInfoContext(ctx, channelID, false)
	if err != nil {
		log.Logger.Error("failed to get channel info from provider", zap.Error(err))
	}

	return chatmodels.Channel{
		ID:   slackChannel.ID,
		Name: slackChannel.Name,
	}, err

}

// DeploymentNotificationThread notifies the thread of the deployment results
func (sp Provider) DeploymentNotificationThread(ctx context.Context, msg, user, channel, ts string) {
	_, _, err := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(userDeploymentNotificationMessage(user, msg), false), slack.MsgOptionTS(ts))
	sp.handleDevOpsErrorNotification(ctx, err)
}

// UserNotificationThread notifies the user in a threaded message
func (sp Provider) UserNotificationThread(ctx context.Context, msg, user, channel, ts string) {
	_, _, err := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(userNotificationMessage(user, msg), false), slack.MsgOptionTS(ts))
	sp.handleDevOpsErrorNotification(ctx, err)
}

// ErrorNotification is a general error notification
func (sp Provider) ErrorNotification(ctx context.Context, user, channel string, err error) {
	log.Logger.Error("slack error notification", zap.Error(err))
	var msg string
	if len(user) > 0 {
		msg = userErrMessage(user, err)
	} else {
		msg = errMessage(err)
	}
	_, _, nerr := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false))
	sp.handleDevOpsErrorNotification(ctx, nerr)
}

// ErrorNotificationThread is a threaded error notification
func (sp Provider) ErrorNotificationThread(ctx context.Context, user, channel, ts string, err error) {
	log.Logger.Error("slack error notification thread", zap.Error(err))
	var msg string
	if len(user) > 0 {
		msg = userErrMessage(user, err)
	} else {
		msg = errMessage(err)
	}
	_, _, nerr := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false), slack.MsgOptionTS(ts))
	sp.handleDevOpsErrorNotification(ctx, nerr)
}

// PostMessageThread sends a threaded message
func (sp Provider) PostMessageThread(ctx context.Context, msg, channel, ts string) (timestamp string) {
	_, respTimestamp, err := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false), slack.MsgOptionTS(ts))
	sp.handleDevOpsErrorNotification(ctx, err)
	return respTimestamp
}

// PostMessage sends a chat message
func (sp Provider) PostMessage(ctx context.Context, msg, channel string) (timestamp string) {
	_, respTS, err := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false))
	sp.handleDevOpsErrorNotification(ctx, err)
	return respTS
}

// GetUser returns user info
func (sp Provider) GetUser(ctx context.Context, user string) (*chatmodels.ChatUser, error) {
	slackUser, err := sp.client.GetUserInfoContext(ctx, user)
	sp.handleDevOpsErrorNotification(ctx, err)
	if err != nil {
		return nil, err
	}
	return mapSlackUser(slackUser), nil
}

// PostLinkMessageThread sends a threaded message with links
func (sp Provider) PostLinkMessageThread(ctx context.Context, url string, user string, channel string, ts string) {

	msgOptionBlocks := slack.MsgOptionBlocks(
		sectionBlockOpt(fmt.Sprintf("<@%s>! %s", user, msgLogLinks)),
		slack.NewDividerBlock(),
		sectionBlockOpt(fmt.Sprintf("<%s|Grafana Logs>", url)),
	)

	linkOpt := slack.MsgOptionEnableLinkUnfurl()
	threadOpt := slack.MsgOptionTS(ts)
	_, _, err := sp.client.PostMessageContext(ctx, channel, msgOptionBlocks, linkOpt, threadOpt)
	sp.handleDevOpsErrorNotification(ctx, err)
}

// ShowResultsMessageThread sends a threaded results message
func (sp Provider) ShowResultsMessageThread(ctx context.Context, msg, user, channel, ts string) {
	msgOptionBlocks := slack.MsgOptionBlocks(
		sectionBlockOpt(fmt.Sprintf("<@%s>! %s", user, msgResultsNotification)),
		slack.NewDividerBlock(),
		sectionBlockOpt(msg),
	)
	threadOpt := slack.MsgOptionTS(ts)
	_, _, err := sp.client.PostMessageContext(ctx, channel, msgOptionBlocks, threadOpt)
	sp.handleDevOpsErrorNotification(ctx, err)
}

// ReleaseResultsMessageThread sens the release results as a threaded message
func (sp Provider) ReleaseResultsMessageThread(ctx context.Context, msg, user, channel, ts string) {
	msgOptionBlocks := slack.MsgOptionBlocks(
		sectionBlockOpt(fmt.Sprintf("<@%s>! %s", user, msgReleaseNotification)),
		slack.NewDividerBlock(),
		sectionBlockOpt(msg),
	)
	threadOpt := slack.MsgOptionTS(ts)
	_, _, err := sp.client.PostMessageContext(ctx, channel, msgOptionBlocks, threadOpt)
	sp.handleDevOpsErrorNotification(ctx, err)
}

func sectionBlockOpt(msg string) *slack.SectionBlock {
	return slack.NewSectionBlock(&slack.TextBlockObject{
		Type:     slack.MarkdownType,
		Text:     msg,
		Emoji:    false,
		Verbatim: false,
	}, nil, nil)
}
