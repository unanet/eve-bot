package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack/slackevents"

	goerror "errors"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// Controller for slack routes
type SlackController struct {
	slackProvider *islack.Provider
}

// New creates a new slack controller (route handler)
func NewSlackController(slackProvider *islack.Provider) *SlackController {
	return &SlackController{
		slackProvider: slackProvider,
	}
}

// Setup the routes
func (c SlackController) Setup(r chi.Router) {
	r.Post("/slack-events", c.slackEventHandler)
	r.Post("/slack-interactive", c.slackInteractiveHandler)
	r.Post("/eve-callback", c.eveCallbackHandler)
	r.Post("/eve-cron-callback", c.eveCronCallbackHandler)

}

func (c SlackController) eveCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the URL Params
	channel := r.URL.Query().Get("channel")
	user := r.URL.Query().Get("user")
	ts := r.URL.Query().Get("ts")
	action := r.URL.Query().Get("action")

	// Get the Body
	payload := eve.NSDeploymentPlan{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.slackProvider.ErrorNotification(r.Context(), user, channel, ts, err)
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	c.slackProvider.EveCallbackNotification(r.Context(), eveapi.CallbackState{User: user, Channel: channel, Payload: payload, TS: ts, Action: action})
	// Just returning an empty response here...
	render.Respond(w, r, nil)
	return

}

func (c SlackController) eveCronCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the URL Params
	channel := r.URL.Query().Get("channel")

	// Get the Body
	payload := eve.NSDeploymentPlan{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.slackProvider.ErrorNotification(r.Context(), "", channel, "", err)
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	user := ""
	if payload.Status == eve.DeploymentPlanStatusErrors {
		user = "channel"
	}

	c.slackProvider.EveCronCallbackNotification(r.Context(), eveapi.CallbackState{User: user, Channel: channel, Payload: payload})
	// Just returning an empty response here...
	render.Respond(w, r, nil)
	return

}

func (c SlackController) slackInteractiveHandler(w http.ResponseWriter, r *http.Request) {
	if err := c.slackProvider.HandleSlackInteraction(r); err != nil {
		render.Respond(w, r, errors.Wrap(err))
		return
	}
	// Just returning an empty response here...
	render.Respond(w, r, nil)
	return
}

func (c SlackController) slackEventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := validateSlackRequest(r, c.slackProvider.Cfg.SlackSigningSecret)
	if err != nil {
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	slackAPIEvent, err := parseSlackEvent(c.slackProvider.Cfg.SlackVerificationToken, body)
	if err != nil {
		render.Respond(w, r, errors.Wrap(botError(err, "failed parse slack event", http.StatusNotAcceptable)))
		return
	}

	// This is a "special" event and only used when setting up the bot endpoint
	if slackAPIEvent.Type == slackevents.URLVerification {
		var slEvent *slackevents.ChallengeResponse
		err := json.Unmarshal(body, &slEvent)
		if err != nil {
			render.Respond(w, r, errors.Wrap(botError(err, "failed to unmarshal slack register event", http.StatusBadRequest)))
			return
		}
		if slEvent == nil || len(slEvent.Challenge) == 0 {
			render.Respond(w, r, errors.Wrap(botError(goerror.New("invalid slack ChallengeResponse event"), "invalid challenge", http.StatusBadGateway)))
			return
		}
		render.Respond(w, r, slEvent.Challenge)
		return
	}

	if err := c.slackProvider.HandleSlackAPIEvent(r.Context(), slackAPIEvent); err != nil {
		render.Respond(w, r, errors.Wrap(err))
		return
	}
	render.Respond(w, r, "OK")
}

func parseSlackEvent(vToken string, body []byte) (slackevents.EventsAPIEvent, error) {
	tokenComp := &slackevents.TokenComparator{VerificationToken: vToken}
	return slackevents.ParseEvent(body, slackevents.OptionVerifyToken(tokenComp))
}

func botError(oerr error, msg string, status int) error {
	log.Logger.Debug("EveBot Error", zap.Error(oerr))
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
