package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/slack-go/slack/slackevents"

	goerror "errors"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// SlackController for slack routes
type SlackController struct {
	svc *service.Provider
}

// NewSlackController creates a new slack controller (route handler)
func NewSlackController(svc *service.Provider) *SlackController {
	return &SlackController{svc: svc}
}

// Setup the routes
func (c SlackController) Setup(r chi.Router) {
	r.Post("/slack-events", c.slackEventHandler)
	r.Post("/slack-interactive", c.slackInteractiveHandler)
}

func (c SlackController) slackInteractiveHandler(w http.ResponseWriter, r *http.Request) {
	if err := c.svc.HandleSlackInteraction(r); err != nil {
		render.Respond(w, r, errors.Wrap(err))
		return
	}
	// Just returning an empty response here...
	render.Respond(w, r, nil)
}

func (c SlackController) slackEventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := validateSlackRequest(r, c.svc.Cfg.SlackSigningSecret)
	if err != nil {
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	slackAPIEvent, err := parseSlackEvent(c.svc.Cfg.SlackVerificationToken, body)
	if err != nil {
		render.Respond(w, r, errors.Wrap(botError(err, "failed parse slack event", http.StatusNotAcceptable)))
		return
	}

	// This is a "special" event and only used when setting up the slackbot endpoint
	if slackAPIEvent.Type == slackevents.URLVerification {
		var slEvent *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &slEvent); err != nil {
			render.Respond(w, r, errors.Wrap(botError(err, "failed to unmarshal slack register event", http.StatusBadRequest)))
			return
		}
		if slEvent == nil || len(slEvent.Challenge) == 0 {
			render.Respond(w, r, errors.Wrap(invalidSlackChallengeError()))
			return
		}
		render.Respond(w, r, slEvent.Challenge)
		return
	}

	log.Logger.Debug("slack event type", zap.Any("event", slackAPIEvent.Type))

	// We are only handling/listening to the CallbackEvent
	if slackAPIEvent.Type != slackevents.CallbackEvent {
		render.Respond(w, r, errors.Wrap(fmt.Errorf("unknown slack API event: %s", slackAPIEvent.Type)))
		return
	}
	innerEvent := slackAPIEvent.InnerEvent
	switch ev := innerEvent.Data.(type) {
	case *slack.FileSharedEvent:
		log.Logger.Info("File Uploaded", zap.Any("event", ev))
	case *slackevents.AppMentionEvent:
		c.svc.HandleSlackAppMentionEvent(r.Context(), ev)
	default:
		log.Logger.Info("slack innerEvent", zap.Any("event", innerEvent))
		render.Respond(w, r, errors.Wrap(unknownSlackEventError(innerEvent)))
		return
	}
	render.Respond(w, r, "OK")
}

func invalidSlackChallengeError() error {
	return botError(
		goerror.New("invalid slack ChallengeResponse event"),
		"invalid challenge",
		http.StatusBadGateway,
	)
}

func unknownSlackEventError(innerEvent slackevents.EventsAPIInnerEvent) error {
	return botError(
		fmt.Errorf("unknown slack inner event: %s", reflect.TypeOf(innerEvent.Data)),
		"unknown slack event",
		http.StatusNotAcceptable,
	)
}

func parseSlackEvent(vToken string, body []byte) (slackevents.EventsAPIEvent, error) {
	tokenComp := &slackevents.TokenComparator{VerificationToken: vToken}
	return slackevents.ParseEvent(body, slackevents.OptionVerifyToken(tokenComp))
}

func botError(oerr error, msg string, status int) error {
	return &errors.RestError{Code: status, Message: msg, OriginalError: oerr}
}

func validateSlackRequest(req *http.Request, signingSecret string) ([]byte, error) {
	verifier, err := slack.NewSecretsVerifier(req.Header, signingSecret)
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
