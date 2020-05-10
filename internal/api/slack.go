package api

import (
	"encoding/json"
	"net/http"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

	"gitlab.unanet.io/devops/eve/pkg/eve"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/errors"
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
	// Payload here is only used for initial URL route verification
	payload, err := c.slackProvider.HandleSlackEvent(r)
	if err != nil {
		render.Respond(w, r, errors.Wrap(err))
		return
	}
	// returning payload response here
	// this is usually just a nil response except for URL verification event
	render.Respond(w, r, payload)
	return
}
