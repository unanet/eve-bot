package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/commander"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// Config needed for slack
//		EVEBOT_SLACK_SIGNING_SECRET
//		EVEBOT_SLACK_VERIFICATION_TOKEN
//		EVEBOT_SLACK_USER_OAUTH_ACCESS_TOKEN
//		EVEBOT_SLACK_OAUTH_ACCESS_TOKEN
//		EVEBOT_SLACK_SKIP_VERIFICATION
type Config struct {
	SlackSigningSecret        string `split_words:"true" required:"true"`
	SlackVerificationToken    string `split_words:"true" required:"true"`
	SlackUserOauthAccessToken string `split_words:"true" required:"true"`
	SlackOauthAccessToken     string `split_words:"true" required:"true"`
	// SlackSkipVerification is used for local dev
	// We need to skip the URL verification when proxying calls with SSH tunnel
	SlackSkipVerification bool `split_words:"true" required:"false"`
}

// Provider provides access to the Slack Client
type Provider struct {
	Client          *slack.Client
	CommandResolver commander.Resolver
	cfg             Config
}

// NewProvider creates a new provider
func NewProvider(cfg Config) *Provider {
	return &Provider{
		Client:          slack.New(cfg.SlackUserOauthAccessToken),
		cfg:             cfg,
		CommandResolver: commander.NewResolver(),
	}
}

func botError(oerr error, msg string, status int) error {
	log.Logger.Error("EveBot Error", zap.Error(oerr))
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
func (p *Provider) HandleInteraction(req *http.Request) error {
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

	// This is required to run the bot locally and proxy calls in with SSH tunnel
	if p.cfg.SlackSkipVerification == false {
		err = verifier.Ensure()
		if err != nil {
			return []byte{}, botError(err, "failed verifier ensure", http.StatusUnauthorized)
		}
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

func (p *Provider) processSlackMentionEvent(ev *slackevents.AppMentionEvent) {
	msgFields := strings.Fields(ev.Text)
	//botIDField := msgFields[0]
	commandFields := msgFields[1:]

	p.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("<@%s>...", ev.User), false))

	eveBotCmd, err := p.CommandResolver.Resolve(commandFields)
	if err != nil {
		p.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Sorry! I can't execute `%s`. Try running `help`...", commandFields), false))
		return
	}

	if eveBotCmd.IsHelpRequest() {
		p.Client.PostMessage(ev.Channel, slack.MsgOptionText(eveBotCmd.Help().String(), false))
		return
	}

	_, err = eveBotCmd.AdditionalArgs()

	if err != nil {
		p.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Whoops! `%s`", err.Error()), false))
		return
	}

	//for _, v := range additionalArgs {
	//	p.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("here is arg key `%s` and value `%v`", v.Name(), v), false))
	//}

	p.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Sure! I'll `%s` that for you. Be right back!", eveBotCmd.Name()), false))
	return
}

// HandleEvent takes an http request and handles the Slack API Event
// this is where we do our request signature validation
// ..and switch the incoming api event types
func (p *Provider) HandleEvent(req *http.Request) (interface{}, error) {
	body, err := p.validateSlackRequest(req)
	if err != nil {
		log.Logger.Debug("Validate Slack Request Error", zap.Error(err))
		return nil, err
	}

	slackAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: p.cfg.SlackVerificationToken}))
	if err != nil {
		return nil, botError(err, "failed parse slack event", http.StatusNotAcceptable)
	}

	log.Logger.Debug("Slack Event Type", zap.String("slack_event", slackAPIEvent.Type))
	switch slackAPIEvent.Type {
	case slackevents.URLVerification:
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &r)
		if err != nil {
			return nil, botError(err, "failed to unmarshal slack reg event", http.StatusBadRequest)
		}
		log.Logger.Debug("Slack Challenge", zap.String("challenge", r.Challenge))
		return r.Challenge, nil
	case slackevents.CallbackEvent:
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			p.processSlackMentionEvent(ev)
			return "OK", nil
		}
	default:
		return nil, fmt.Errorf("unknown slack event: %v", slackAPIEvent.Type)
	}
	return nil, fmt.Errorf("unknown slack event: %v", slackAPIEvent.Type)
}
