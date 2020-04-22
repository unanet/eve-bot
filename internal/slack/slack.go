package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve/pkg/eveerrs"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// Config needed for slack
type Config struct {
	SlackSigningSecret     string `split_words:"true" required:"true"`
	SlackVerificationToken string `split_words:"true" required:"true"`
	SlackbotOauthToken     string `split_words:"true" required:"true"`
	SlackOauthToken        string `split_words:"true" required:"true"`
	// SlackSkipVerification is used for local dev
	// We need to skip the URL verification when proxying calls with SSH tunnel
	SlackSkipVerification bool `split_words:"true" required:"false"`
}

// Provider provides access to the Slack Client
type Provider struct {
	Client *slack.Client
	cfg    Config
}

// NewProvider creates a new provider
func NewProvider(cfg Config) *Provider {
	return &Provider{
		Client: slack.New(cfg.SlackbotOauthToken),
		cfg:    cfg,
	}
}

func botError(oerr error, msg string, status int) error {
	log.Logger.Error("evebot error", zap.Error(oerr))
	return &eveerrs.RestError{
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
func (p *Provider) HandleInteraction(req *http.Request) error {
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
func (p *Provider) HandleEvent(req *http.Request) error {
	verifier, err := slack.NewSecretsVerifier(req.Header, p.cfg.SlackSigningSecret)
	if err != nil {
		return botError(err, "failed new secret verifier", http.StatusUnauthorized)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return botError(err, "failed readAll req body", http.StatusBadRequest)
	}

	_, err = verifier.Write(body)
	if err != nil {
		return botError(err, "failed verifier write", http.StatusUnauthorized)
	}

	// This is required to run the bot locally and proxy calls in with SSH tunnel
	if p.cfg.SlackSkipVerification == false {
		err = verifier.Ensure()
		if err != nil {
			return botError(err, "failed verifier ensure", http.StatusUnauthorized)
		}
	}

	slackAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: p.cfg.SlackVerificationToken}))
	if err != nil {
		return botError(err, "failed parse slack event", http.StatusNotAcceptable)
	}

	log.Logger.Debug("slack event", zap.String("slack_event_type", slackAPIEvent.Type))
	switch slackAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return botError(err, "failed to unmarshal slack reg event", http.StatusBadRequest)
		}
	case slackevents.CallbackEvent:
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			msgFields := strings.Fields(ev.Text)
			botIDField := msgFields[0]
			commandFields := msgFields[1:]

			if len(commandFields) <= 0 || commandFields[0] == "help" {
				return p.ShowHelp(ev)
			}

			p.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Yes, <@%s>, this is %s! You want me to run: %s", ev.User, botIDField, commandFields), false))
		}
	default:
		return fmt.Errorf("unknown slack event: %v", slackAPIEvent.Type)
	}

	return nil
}

// ShowHelp shows the help message to the Slack User
func (p *Provider) ShowHelp(ev *slackevents.AppMentionEvent) error {

	helpAttachment := slack.Attachment{
		Pretext:    "\ndeploy help\nmigrate help",
		Fallback:   "help",
		CallbackID: "help",
		Color:      "#3AA3E3",
	}
	attachmentOpt := slack.MsgOptionAttachments(helpAttachment)
	msgOpt := slack.MsgOptionText(fmt.Sprintf("Hey <@%s>! Need a little help? Try the following...", ev.User), false)
	p.Client.PostMessage(ev.Channel, msgOpt, attachmentOpt)
	return nil
}
