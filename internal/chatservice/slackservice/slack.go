package slackservice

import (
	"context"

	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/chatmodels"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

type Provider struct {
	client *slack.Client
}

func New(c *slack.Client) Provider {
	return Provider{client: c}
}

func (sp Provider) handleDevOpsErrorNotification(err error) {
	if err != nil {
		log.Logger.Error("critical devops error", zap.Error(err))
		sp.client.PostMessage(devOpsMonitoringChannel, slack.MsgOptionText(errMessage(err), false))
	}
}

func (sp Provider) DeploymentNotificationThread(ctx context.Context, msg, user, channel, ts string) {
	log.Logger.Debug("deployment notification", zap.String("user", user), zap.String("message", msg))
	_, _, err := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(userDeploymentNotificationMessage(user, msg), false), slack.MsgOptionTS(ts))
	sp.handleDevOpsErrorNotification(err)
}

func (sp Provider) UserNotificationThread(ctx context.Context, msg, user, channel, ts string) {
	log.Logger.Debug("user notification", zap.String("user", user), zap.String("message", msg))
	_, _, err := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(userNotificationMessage(user, msg), false), slack.MsgOptionTS(ts))
	sp.handleDevOpsErrorNotification(err)
}

func (sp Provider) ErrorNotification(ctx context.Context, user, channel string, err error) {
	log.Logger.Error("slack error notification", zap.Error(err))
	var msg string
	if len(user) > 0 {
		msg = userErrMessage(user, err)
	} else {
		msg = errMessage(err)
	}
	_, _, nerr := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false))
	sp.handleDevOpsErrorNotification(nerr)
}

func (sp Provider) ErrorNotificationThread(ctx context.Context, user, channel, ts string, err error) {
	log.Logger.Error("slack error notification thread", zap.Error(err))
	var msg string
	if len(user) > 0 {
		msg = userErrMessage(user, err)
	} else {
		msg = errMessage(err)
	}
	_, _, nerr := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false), slack.MsgOptionTS(ts))
	sp.handleDevOpsErrorNotification(nerr)
}

func (sp Provider) PostMessageThread(ctx context.Context, msg, channel, ts string) (timestamp string) {
	_, respTimestamp, err := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false), slack.MsgOptionTS(ts))
	sp.handleDevOpsErrorNotification(err)
	return respTimestamp
}

func (sp Provider) PostMessage(ctx context.Context, msg, channel string) {
	_, _, err := sp.client.PostMessageContext(ctx, channel, slack.MsgOptionText(msg, false))
	sp.handleDevOpsErrorNotification(err)
}

func (sp Provider) GetUser(ctx context.Context, user string) (*chatmodels.ChatUser, error) {
	slackUser, err := sp.client.GetUserInfoContext(ctx, user)
	sp.handleDevOpsErrorNotification(err)
	if err != nil {
		return nil, err
	}
	return mapSlackUser(slackUser), nil
}
