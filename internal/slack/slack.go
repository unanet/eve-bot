package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

func (p *Provider) ErrorNotification(ctx context.Context, user, channel string, err error) error {
	slackErrMsg := fmt.Sprintf("Sorry <@%s>! Something terrible has happened:\n\n ```%v```\n\nI've dispatched a semi-competent team of monkeys to battle the issue...", user, err.Error())
	_, _, _ = p.Client.PostMessageContext(
		ctx,
		channel,
		slack.MsgOptionText(slackErrMsg, false))
	return nil
}

// HandleEveCallback handles the callbacks from eve-api
func (p *Provider) EveCallbackNotification(ctx context.Context, cbState eveapi.CallbackState) error {

	headerOpt := slack.MsgOptionText(cbState.SlackMsgHeader(), false)

	resultOpt := newBlockMsgOpt(cbState.SlackMsgResults())

	_, _, _ = p.Client.PostMessageContext(ctx, cbState.Channel, headerOpt, resultOpt)
	return nil
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

func newBlockMsgOpt(text string) slack.MsgOption {
	return slack.MsgOptionBlocks(
		slack.NewSectionBlock(
			slack.NewTextBlockObject(
				slack.MarkdownType,
				text,
				true,
				false),
			nil,
			nil),
		slack.NewDividerBlock())
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
			p.Client.PostMessageContext(req.Context(), ev.Channel, slack.MsgOptionText(cmd.AckMsg(ev.User), false))

			if cmd.MakeAsyncReq() {
				// Call API in separate Go Routine
				go func() {
					apiReqObj := cmd.EveReqObj(callBackURL)

					switch apiReqObj.(type) {
					case eveapi.DeploymentPlanOptions:
						resp, err := p.EveAPIClient.Deploy(context.TODO(), apiReqObj.(eveapi.DeploymentPlanOptions), ev.User, ev.Channel)
						if err != nil {
							log.Logger.Debug("eve-api error", zap.Error(err))

							p.Client.PostMessageContext(
								context.TODO(),
								ev.Channel,
								slack.MsgOptionText(
									fmt.Sprintf("Whoops <@%s>! I detected some *errors:*\n\n ```%v```", ev.User, err.Error()), false))
							return
						}
						log.Logger.Debug("eve-api response", zap.Any("response", resp))
					default:
						log.Logger.Error("invalid eve api command request object")
					}

				}()
			}

			// Immediately respond to the Slack HTTP Request.
			// This doesn't actually do anything except free up the incoming request
			return "OK", nil
		}
	default:
		return nil, fmt.Errorf("unknown slack event: %v", slackAPIEvent.Type)
	}
	return nil, fmt.Errorf("unknown slack event: %v", slackAPIEvent.Type)
}

//	queue.WorkQueue <- queue.WorkRequest{
//		Name:    ev.Channel,
//		User:    ev.User,
//		Channel: ev.Channel,
//		EveType: cmd.Name(),
//		Delay:   time.Second * 60, // Just for testing/simulation
//	}
