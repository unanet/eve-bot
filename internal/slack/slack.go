package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve/pkg/eveerrs"
)

// Config needed for slack
type Config struct {
	SlackSigningSecret     string `split_words:"true" required:"true"`
	SlackVerificationToken string `split_words:"true" required:"true"`
	SlackBotOAuthToken     string `split_words:"true" required:"true"`
	SlackOAuthToken        string `split_words:"true" required:"true"`
}

// Provider provides access to the Slack Client
// basically a wrapper around slack
type Provider interface {
	HandleEvent(req *http.Request) error
}

type provider struct {
	client *slack.Client
	cfg    Config
}

// NewProvider creates a new provider
func NewProvider(cfg Config) Provider {
	return &provider{
		client: slack.New(cfg.SlackBotOAuthToken),
		cfg:    cfg,
	}
}

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *provider) HandleEvent(req *http.Request) error {

	restErr := func(oerr error, msg string, status int) error {
		return &eveerrs.RestError{
			Code:          http.StatusUnauthorized,
			Message:       msg,
			OriginalError: oerr,
		}
	}

	verifier, err := slack.NewSecretsVerifier(req.Header, p.cfg.SlackSigningSecret)
	if err != nil {
		return restErr(err, "failed new secret verifier", http.StatusUnauthorized)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return restErr(err, "failed readAll req body", http.StatusBadRequest)
	}

	_, err = verifier.Write(body)
	if err != nil {
		return restErr(err, "failed verifier write", http.StatusUnauthorized)
	}

	err = verifier.Ensure()
	if err != nil {
		return restErr(err, "failed verifier ensure", http.StatusUnauthorized)
	}

	slackAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: p.cfg.SlackVerificationToken}))

	if err != nil {
		return restErr(err, "failed parse slack event", http.StatusNotAcceptable)
	}

	// p.logger.For(req.Context()).Info("slack event", zap.String("event", slackAPIEvent.Type))

	if slackAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return restErr(err, "failed to unmarshal slack reg event", http.StatusBadRequest)
		}
	}

	if slackAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			p.client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, App Mention Event hello.", false))
		}
	}

	return nil
}
