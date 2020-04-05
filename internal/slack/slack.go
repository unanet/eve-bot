package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack/slackevents"
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/api/resterror"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/evelogger"
	"go.uber.org/zap"
)

// Provider provides access to the Slack Client
// basically a wrapper around slack
type Provider interface {
	HandleEvent(req *http.Request) error
}

type provider struct {
	client *slack.Client
	cfg    *config.Config
	logger evelogger.Container
}

// NewProvider creates a new provider
func NewProvider(cfg *config.Config, logger evelogger.Container) Provider {
	return &provider{
		client: slack.New(cfg.SlackSecrets.BotOAuthToken),
		cfg:    cfg,
		logger: logger,
	}
}

func (p *provider) Client() *slack.Client {
	return p.client
}

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *provider) HandleEvent(req *http.Request) error {

	restErr := func(oerr error, msg string, status int) error {
		return &resterror.RestError{
			Code:          http.StatusUnauthorized,
			Message:       msg,
			OriginalError: oerr,
		}
	}

	verifier, err := slack.NewSecretsVerifier(req.Header, p.cfg.SlackSecrets.SigningSecret)
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
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: p.cfg.SlackSecrets.VerificationToken}))

	if err != nil {
		return restErr(err, "failed parse slack event", http.StatusNotAcceptable)
	}

	p.logger.For(req.Context()).Info("slack event", zap.String("event", slackAPIEvent.Type))

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
