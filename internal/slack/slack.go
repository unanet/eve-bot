package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	eveerrs "gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// EveCallbackNotification handles the callbacks from eve-api and notifies the slack user
func (p *Provider) EveCallbackNotification(ctx context.Context, cbState eveapi.CallbackState) {
	log.Logger.Debug("eve callback notification", zap.Any("cb_state", cbState))
	respChannel, respTS, err := p.Client.PostMessageContext(ctx, cbState.Channel, slack.MsgOptionText(cbState.ToChatMsg(), false), slack.MsgOptionTS(cbState.TS))
	if err != nil {
		p.ErrorNotification(ctx, cbState.User, respChannel, respTS, err)
	}
}

func (p *Provider) EveCronCallbackNotification(ctx context.Context, cbState eveapi.CallbackState) {
	log.Logger.Debug("eve cron callback notification", zap.Any("cb_state", cbState))
	// ignore pending messages
	if cbState.Payload.Status == eve.DeploymentPlanStatusPending {
		log.Logger.Debug("eve cron callback notification pending", zap.Any("status", cbState.Payload.Status))
		return
	}

	respChannel, respTS, err := p.Client.PostMessageContext(ctx, cbState.Channel, slack.MsgOptionText(cbState.ToChatMsg(), false))
	if err != nil {
		p.ErrorNotification(ctx, cbState.User, respChannel, respTS, err)
	}
}

func (p *Provider) handleEveApiResponse(slackUser, slackChannel, respTS string, resp *eveapi.DeploymentPlanOptions, err error) {
	// The err is coming back with an empty message...
	if err != nil && len(err.Error()) > 0 {
		log.Logger.Debug("deployment error notification", zap.Error(err))
		p.Client.PostMessageContext(
			context.TODO(),
			slackChannel,
			slack.MsgOptionText(fmt.Sprintf("<@%s>! %s\n\n ```%s```", slackUser, msgDeploymentErrNotification, err.Error()), false),
			slack.MsgOptionTS(respTS),
		)
		return
	}

	if resp == nil {
		p.ErrorNotification(context.TODO(), slackUser, slackChannel, respTS, errors.New("invalid api response"))
		return
	}

	if len(resp.Messages) > 0 {
		msg := strings.Join(resp.Messages, ",")
		log.Logger.Debug("slack deployment message notification", zap.String("message", msg))
		p.Client.PostMessageContext(
			context.TODO(),
			slackChannel,
			slack.MsgOptionText(fmt.Sprintf("<@%s>! %s\n\n ```%s```\n\n", slackUser, msgNotification, msg), false),
			slack.MsgOptionTS(respTS),
		)
		return
	}
}

// HandleInteraction handles the interactive callbacks (buttons, dropdowns, etc.)
func (p *Provider) HandleSlackInteraction(req *http.Request) error {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(req.FormValue("payload")), &payload)
	if err != nil {
		return botError(err, "failed to parse interactive slack message payload", http.StatusInternalServerError)
	}
	log.Logger.Info(fmt.Sprintf("Message button pressed by user %s with value %s", payload.User.Name, payload.Value))
	return nil
}

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *Provider) HandleSlackAPIEvent(ctx context.Context, slackAPIEvent slackevents.EventsAPIEvent) error {
	switch slackAPIEvent.Type {
	case slackevents.CallbackEvent:
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// Resolve the input and return a Command object
			cmd := p.CommandResolver.Resolve(ev.Text, ev.Channel, ev.User, ev.ThreadTimeStamp)
			// Send the immediate Acknowledgement Message back to the chat user
			_, respTS, _ := p.Client.PostMessageContext(ctx, ev.Channel, slack.MsgOptionText(cmd.AckMsg(ev.User), false), slack.MsgOptionTS(ev.ThreadTimeStamp))
			realUser, err := p.Client.GetUserInfo(ev.User)
			if err != nil {
				p.ErrorNotification(context.TODO(), ev.User, ev.Channel, respTS, fmt.Errorf("failed to get user details"))
				return nil
			}
			//go cmd.Execute()
			//// Immediately respond to the Slack HTTP Request.
			//return "OK", nil

			if cmd.MakeAsyncReq() {
				// Call API in separate Go Routine
				go func(reqObj interface{}, slackUser, slackChannel, respTS string) {
					if reqObj == nil {
						p.ErrorNotification(context.TODO(), slackUser, slackChannel, respTS, errInvalidRequestObj)
						return
					}
					switch reqObj.(type) {
					case eveapi.DeploymentPlanOptions:
						resp, err := p.EveAPIClient.Deploy(context.TODO(), reqObj.(eveapi.DeploymentPlanOptions), slackUser, slackChannel, respTS)
						p.handleEveApiResponse(slackUser, slackChannel, respTS, resp, err)
						return
					default:
						p.ErrorNotification(context.TODO(), slackUser, slackChannel, respTS, errInvalidRequestObj)
						return
					}
				}(cmd.EveReqObj(callBackURL, realUser.Name), ev.User, ev.Channel, respTS)
			}
			//Immediately respond to the Slack HTTP Request.
			return nil
		default:
			return &eveerrs.RestError{
				Code:          http.StatusBadRequest,
				Message:       "unknown slack event",
				OriginalError: fmt.Errorf("unknown slack inner event: %s", reflect.TypeOf(innerEvent.Data)),
			}
		}
	}
	return fmt.Errorf("unknown slack API event: %s", slackAPIEvent.Type)
}
