package slack

import (
	"context"
	"errors"
	"fmt"

	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

var (
	errInvalidRequestObj = errors.New("invalid request object")
)

const (
	msgErrNotification           = "Something terrible has happened..."
	msgErrNotificationAssurance  = "I've dispatched a semi-competent team of monkeys to battle the issue..."
	msgNotification              = "I've got some news..."
	msgDeploymentErrNotification = "I detected some deployment *errors:*"
)

// ErrorNotification is generic error function so that we can message to the slack user that something bad has happened
// we should probably have
func (p *Provider) ErrorNotification(ctx context.Context, user, channel, ts string, err error) {
	log.Logger.Error("slack error notification", zap.Error(err))
	var slackErrMsg string
	if len(user) > 0 {
		slackErrMsg = fmt.Sprintf("<@%s>! %s\n\n ```%s```\n\n%s", user, msgErrNotification, err.Error(), msgErrNotificationAssurance)
	} else {
		slackErrMsg = fmt.Sprintf("%s\n\n ```%s```\n\n%s", msgErrNotification, err.Error(), msgErrNotificationAssurance)
	}

	if len(ts) > 0 {
		_, _, _ = p.Client.PostMessageContext(ctx, channel, slack.MsgOptionText(slackErrMsg, false), slack.MsgOptionTS(ts))
	} else {
		_, _, _ = p.Client.PostMessageContext(ctx, channel, slack.MsgOptionText(slackErrMsg, false))
	}
}
