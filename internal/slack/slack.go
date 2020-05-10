package slack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// ErrorNotification is generic error function so that we can message to the slack user that something bad has happened
// we should probably have
func (p *Provider) ErrorNotification(ctx context.Context, user, channel string, err error) {
	log.Logger.Error("slack error notification", zap.Error(err))
	slackErrMsg := fmt.Sprintf("Sorry <@%s>! Something terrible has happened:\n\n ```%s```\n\nI've dispatched a semi-competent team of monkeys to battle the issue...", user, err.Error())
	_, _, _ = p.Client.PostMessageContext(ctx, channel, slack.MsgOptionText(slackErrMsg, false))
}

// EveCallbackNotification handles the callbacks from eve-api and notifies the slack user
func (p *Provider) EveCallbackNotification(ctx context.Context, cbState eveapi.CallbackState) {
	_, _, err := p.Client.PostMessageContext(ctx, cbState.Channel, slack.MsgOptionText(cbState.ToChatMsg(), false))
	if err != nil {
		p.ErrorNotification(ctx, cbState.User, cbState.Channel, err)
	}
}

func (p *Provider) messageNotification(ctx context.Context, user, channel, message string) {
	log.Logger.Debug("slack deployment message notification", zap.String("message", message))
	slackErrMsg := fmt.Sprintf("<@%s>! Something has happened:\n\n ```%s```\n\n", user, message)
	_, _, _ = p.Client.PostMessageContext(ctx, channel, slack.MsgOptionText(slackErrMsg, false))
}

func (p *Provider) deploymentErrorNotification(ctx context.Context, user, channel string, err error) {
	log.Logger.Debug("deployment error notification", zap.Error(err))
	slackDeploymentErrMsg := fmt.Sprintf("Whoops <@%s>! I detected some deployment *errors:*\n\n ```%s```", user, err.Error())
	_, _, _ = p.Client.PostMessageContext(ctx, channel, slack.MsgOptionText(slackDeploymentErrMsg, false))
}

func (p *Provider) handleEveApiResponse(slackUser, slackChannel string, resp *eveapi.DeploymentPlanOptions, err error) {
	if err != nil {
		p.deploymentErrorNotification(context.TODO(), slackUser, slackChannel, err)
		return
	}

	if resp == nil {
		p.ErrorNotification(context.TODO(), slackUser, slackChannel, errInvalidEveApiResponse)
		return
	}

	if len(resp.Messages) > 0 {
		p.messageNotification(context.TODO(), slackUser, slackChannel, strings.Join(resp.Messages, ","))
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

	fmt.Printf("Message button pressed by user %s with value %s", payload.User.Name, payload.Value)
	return nil
}

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *Provider) HandleSlackEvent(req *http.Request) (interface{}, error) {
	body, err := validateSlackRequest(req)
	if err != nil {
		log.Logger.Debug("Validate Slack Request Error", zap.Error(err))
		return nil, err
	}

	slackAPIEvent, err := slackevents.ParseEvent(body,
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{
				VerificationToken: p.cfg.SlackVerificationToken,
			},
		),
	)

	if err != nil {
		return nil, botError(err, "failed parse slack event", http.StatusNotAcceptable)
	}

	switch slackAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			return nil, botError(err, "failed to unmarshal slack reg event", http.StatusBadRequest)
		}
		return r.Challenge, nil
	case slackevents.CallbackEvent:
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// Resolve the input and return a Command object
			cmd := p.CommandResolver.Resolve(ev.Text)
			// Send the immediate Acknowledgement Message back to the chat user
			_, _, _ = p.Client.PostMessageContext(req.Context(), ev.Channel, slack.MsgOptionText(cmd.AckMsg(ev.User), false))

			if cmd.MakeAsyncReq() {
				// Call API in separate Go Routine
				go func(reqObj interface{}, slackUser, slackChannel string) {
					if reqObj == nil {
						p.ErrorNotification(context.TODO(), slackUser, slackChannel, errInvalidRequestObj)
						return
					}
					switch reqObj.(type) {
					case eveapi.DeploymentPlanOptions:
						resp, err := p.EveAPIClient.Deploy(context.TODO(), reqObj.(eveapi.DeploymentPlanOptions), slackUser, slackChannel)
						p.handleEveApiResponse(slackUser, slackChannel, resp, err)

					default:
						p.ErrorNotification(context.TODO(), slackUser, slackChannel, errInvalidRequestObj)
						return
					}
				}(cmd.EveReqObj(callBackURL), ev.User, ev.Channel)
			}
			// Immediately respond to the Slack HTTP Request.
			return "OK", nil
		}
	default:
		return nil, unknownSlackEventErr(slackAPIEvent.Type)
	}
	return nil, unknownSlackEventErr(slackAPIEvent.Type)
}

func unknownSlackEventErr(slackEvent string) error {
	return fmt.Errorf("unknown slack event: %s", slackEvent)
}

var (
	errInvalidRequestObj     = errors.New("invalid request object")
	errInvalidEveApiResponse = errors.New("invalid api response")
)
