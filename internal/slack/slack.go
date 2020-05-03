package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// Config needed for slack
//		EVEBOT_SLACK_SIGNING_SECRET
//		EVEBOT_SLACK_VERIFICATION_TOKEN
//		EVEBOT_SLACK_USER_OAUTH_ACCESS_TOKEN
//		EVEBOT_SLACK_OAUTH_ACCESS_TOKEN
type Config struct {
	SlackSigningSecret        string `split_words:"true" required:"true"`
	SlackVerificationToken    string `split_words:"true" required:"true"`
	SlackUserOauthAccessToken string `split_words:"true" required:"true"`
	SlackOauthAccessToken     string `split_words:"true" required:"true"`
}

// Provider provides access to the Slack Client
type Provider struct {
	Client          *slack.Client
	CommandResolver botcommander.Resolver
	EveAPIClient    eveapi.Client
	cfg             Config
}

var (
	callBackURL string
)

// NewProvider creates a new provider
func NewProvider(sClient *slack.Client, commander botcommander.Resolver, eveAPIClient eveapi.Client, cfg Config) *Provider {
	callBackURL = eveAPIClient.CallBackURL()
	return &Provider{
		Client:          sClient,
		CommandResolver: commander,
		EveAPIClient:    eveAPIClient,
		cfg:             cfg,
	}
}

func botError(oerr error, msg string, status int) error {
	log.Logger.Debug("EveBot Error", zap.Error(oerr))
	return &errors.RestError{
		Code:          status,
		Message:       msg,
		OriginalError: oerr,
	}
}

// HandleEveCallback handles the callbacks from eve-api
func (p *Provider) HandleEveCallback(req *http.Request) error {
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

func (p *Provider) validateSlackRequest(req *http.Request) ([]byte, error) {
	verifier, err := slack.NewSecretsVerifier(req.Header, p.cfg.SlackSigningSecret)
	if err != nil {
		return []byte{}, botError(err, "failed new secret verifier", http.StatusUnauthorized)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return []byte{}, botError(err, "failed readAll req body", http.StatusBadRequest)
	}

	_, err = verifier.Write(body)
	if err != nil {
		return []byte{}, botError(err, "failed verifier write", http.StatusUnauthorized)
	}

	err = verifier.Ensure()
	if err != nil {
		// Sending back a Teapot StatusCode here (418)
		// These are requests from bad actors
		return []byte{}, botError(err, "failed verifier ensure", http.StatusTeapot)
	}

	return body, nil
}

func newBlockMsgOpt(text string) slack.MsgOption {
	return slack.MsgOptionBlocks(
		slack.NewSectionBlock(
			slack.NewTextBlockObject(
				slack.MarkdownType,
				text,
				false,
				false),
			nil,
			nil),
		slack.NewDividerBlock())
}

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *Provider) HandleSlackEvent(req *http.Request) (interface{}, error) {
	body, err := p.validateSlackRequest(req)
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

			// cmd.WorkRequest handles sync vs async
			// Make the work Request to API here:
			//p.Client.PostMessage(ev.Channel, slack.MsgOptionText(cmd.WorkRequest(ev.Channel, ev.User), false))
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
