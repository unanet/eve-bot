package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/eve"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	eveerrs "gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

var (
	errInvalidRequestObj     = errors.New("invalid request object")
	errInvalidEveApiResponse = errors.New("invalid api response")
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

// EveCallbackNotification handles the callbacks from eve-api and notifies the slack user
func (p *Provider) EveCallbackNotification(ctx context.Context, cbState eveapi.CallbackState) {
	respChannel, respTS, err := p.Client.PostMessageContext(ctx, cbState.Channel, slack.MsgOptionText(cbState.ToChatMsg(), false), slack.MsgOptionTS(cbState.TS))
	if err != nil {
		p.ErrorNotification(ctx, cbState.User, respChannel, respTS, err)
	}
}

func (p *Provider) EveCronCallbackNotification(ctx context.Context, cbState eveapi.CallbackState) {
	// ignore pending messages
	if cbState.Payload.Status == eve.DeploymentPlanStatusPending {
		return
	}

	respChannel, respTS, err := p.Client.PostMessageContext(ctx, cbState.Channel, slack.MsgOptionText(cbState.ToChatMsg(), false))
	if err != nil {
		p.ErrorNotification(ctx, cbState.User, respChannel, respTS, err)
	}
}

func (p *Provider) messageNotification(ctx context.Context, user, channel, respTS, message string) {
	log.Logger.Debug("slack deployment message notification", zap.String("message", message))
	slackErrMsg := fmt.Sprintf("<@%s>! %s\n\n ```%s```\n\n", user, msgNotification, message)
	_, _, _ = p.Client.PostMessageContext(ctx, channel, slack.MsgOptionText(slackErrMsg, false), slack.MsgOptionTS(respTS))
}

func (p *Provider) deploymentErrorNotification(ctx context.Context, user, channel, respTS string, err error) {
	log.Logger.Debug("deployment error notification", zap.Error(err))
	slackDeploymentErrMsg := fmt.Sprintf("<@%s>! %s\n\n ```%s```", user, msgDeploymentErrNotification, err.Error())
	_, _, _ = p.Client.PostMessageContext(ctx, channel, slack.MsgOptionText(slackDeploymentErrMsg, false), slack.MsgOptionTS(respTS))
}

func (p *Provider) handleEveApiResponse(slackUser, slackChannel, respTS string, resp *eveapi.DeploymentPlanOptions, err error) {
	// The err is coming back with an empty message...
	if err != nil && len(err.Error()) > 0 {
		p.deploymentErrorNotification(context.TODO(), slackUser, slackChannel, respTS, err)
		return
	}

	if resp == nil {
		p.ErrorNotification(context.TODO(), slackUser, slackChannel, respTS, errInvalidEveApiResponse)
		return
	}

	if len(resp.Messages) > 0 {
		p.messageNotification(context.TODO(), slackUser, slackChannel, respTS, strings.Join(resp.Messages, ","))
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

func parseSlackEvent(vToken string, body []byte) (slackevents.EventsAPIEvent, error) {
	tokenComp := &slackevents.TokenComparator{VerificationToken: vToken}
	tokenCompOpt := slackevents.OptionVerifyToken(tokenComp)
	return slackevents.ParseEvent(body, tokenCompOpt)
}

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *Provider) HandleSlackEvent(ctx context.Context, body []byte) (interface{}, error) {
	slackAPIEvent, err := parseSlackEvent(p.Cfg.SlackVerificationToken, body)
	if err != nil {
		return nil, &eveerrs.RestError{Code: http.StatusNotAcceptable, Message: "failed parse slack event", OriginalError: err}
	}

	switch slackAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			return nil, &eveerrs.RestError{
				Code:          http.StatusBadRequest,
				Message:       "failed to unmarshal slack reg event",
				OriginalError: err,
			}
		}
		return r.Challenge, nil
	case slackevents.CallbackEvent:
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// Resolve the input and return a Command object
			cmd := p.CommandResolver.Resolve(ev.Text)
			// Send the immediate Acknowledgement Message back to the chat user
			_, respTS, _ := p.Client.PostMessageContext(ctx, ev.Channel, slack.MsgOptionText(cmd.AckMsg(ev.User), false), slack.MsgOptionTS(ev.ThreadTimeStamp))
			realUser, err := p.Client.GetUserInfo(ev.User)
			if err != nil {
				p.ErrorNotification(context.TODO(), ev.User, ev.Channel, respTS, fmt.Errorf("failed to get user details"))
				return "OK", nil
			}
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
			// Immediately respond to the Slack HTTP Request.
			return "OK", nil
		default:
			return nil, &eveerrs.RestError{
				Code:          http.StatusBadRequest,
				Message:       "unknown slack event",
				OriginalError: fmt.Errorf("unknown slack inner event: %s", reflect.TypeOf(innerEvent.Data)),
			}
		}
	}
	return nil, fmt.Errorf("unknown slack API event: %s", slackAPIEvent.Type)
}
