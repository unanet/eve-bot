package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
}

// Provider provides access to the Slack Client
// basically a wrapper around slack
// type Provider interface {
// 	HandleEvent(req *http.Request) error
// }

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

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *Provider) HandleEvent(req *http.Request) error {
	log.Logger.Info("SlackClient HandleEvent")

	restErr := func(oerr error, msg string, status int) error {
		log.Logger.Error("SlackClient HandleEvent Error", zap.Error(oerr))
		return &eveerrs.RestError{
			Code:          http.StatusUnauthorized,
			Message:       msg,
			OriginalError: oerr,
		}
	}

	log.Logger.Info("SlackClient NewSecretsVerifier")
	verifier, err := slack.NewSecretsVerifier(req.Header, p.cfg.SlackSigningSecret)
	if err != nil {
		return restErr(err, "failed new secret verifier", http.StatusUnauthorized)
	}

	log.Logger.Info("SlackClient ReadBody")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return restErr(err, "failed readAll req body", http.StatusBadRequest)
	}

	log.Logger.Info("SlackClient verifier.Write(body)")
	_, err = verifier.Write(body)
	if err != nil {
		return restErr(err, "failed verifier write", http.StatusUnauthorized)
	}

	log.Logger.Info("SlackClient verifier.Ensure")
	err = verifier.Ensure()
	if err != nil {
		return restErr(err, "failed verifier ensure", http.StatusUnauthorized)
	}

	log.Logger.Info("SlackClient slackevents.ParseEvent")
	slackAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: p.cfg.SlackVerificationToken}))

	if err != nil {
		return restErr(err, "failed parse slack event", http.StatusNotAcceptable)
	}

	// p.logger.For(req.Context()).Info("slack event", zap.String("event", slackAPIEvent.Type))
	log.Logger.Info("SlackClient slackAPIEvent.Type", zap.String("slack_event_type", slackAPIEvent.Type))
	if slackAPIEvent.Type == slackevents.URLVerification {
		log.Logger.Info("SlackClient slackAPIEvent.Type")
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return restErr(err, "failed to unmarshal slack reg event", http.StatusBadRequest)
		}
	}

	if slackAPIEvent.Type == slackevents.CallbackEvent {
		log.Logger.Info("SlackClient slackAPIEvent.Type")
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			p.Client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, App Mention Event hello.", false))
		}
	}

	return nil
}
