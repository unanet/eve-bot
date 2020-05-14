package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
)

const (
	DevopsSlackGroupID = "S013MD29N3X"
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

	// Get the Body
	payload := eve.NSDeploymentPlan{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.slackProvider.ErrorNotification(r.Context(), user, channel, err)
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	c.slackProvider.EveCallbackNotification(r.Context(), eveapi.CallbackState{User: user, Channel: channel, Payload: payload})
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
		c.slackProvider.ErrorNotification(r.Context(), DevopsSlackGroupID, channel, err)
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	c.slackProvider.EveCallbackNotification(r.Context(), eveapi.CallbackState{User: DevopsSlackGroupID, Channel: channel, Payload: payload})
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
	// Payload here is only used for initial URL route verification
	payload, err := c.slackProvider.HandleSlackEvent(r.Context(), body)
	if err != nil {
		render.Respond(w, r, errors.Wrap(err))
		return
	}
	// returning payload response here
	// this is usually just a nil response except for URL verification event
	render.Respond(w, r, payload)
	return
}

func botError(oerr error, msg string, status int) error {
	log.Logger.Debug("EveBot Error", zap.Error(oerr))
	return &errors.RestError{
		Code:          status,
		Message:       msg,
		OriginalError: oerr,
	}
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
