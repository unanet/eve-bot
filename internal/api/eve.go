package api

import (
	"encoding/json"
	"net/http"

	"github.com/unanet/eve-bot/internal/eveapi"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/unanet/eve-bot/internal/service"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/go/pkg/errors"
)

// EveController for slack routes
type EveController struct {
	svc *service.Provider
}

// NewEveController New creates a new eve controller (route handler)
func NewEveController(svc *service.Provider) *EveController {
	return &EveController{svc: svc}
}

// Setup the routes
func (c EveController) Setup(r chi.Router) {
	r.Post("/eve-callback", c.eveCallbackHandler)
	r.Post("/eve-cron-callback", c.eveCronCallbackHandler)
}

func logLink(ns string) string {
	return "https://grafana.plainsight.biz/"
}

func (c EveController) eveCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the URL Params
	channel := r.URL.Query().Get("channel")
	user := r.URL.Query().Get("user")
	ts := r.URL.Query().Get("ts")

	// Get the Body
	payload := eve.NSDeploymentPlan{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.svc.ChatService.ErrorNotificationThread(r.Context(), user, channel, ts, err)
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	cbState := eveapi.CallbackState{User: user, Channel: channel, Payload: payload, TS: ts}
	c.svc.ChatService.PostMessageThread(r.Context(), cbState.ToChatMsg(), cbState.Channel, cbState.TS)

	if cbState.Payload.Status == eve.DeploymentPlanStatusErrors {
		c.svc.ChatService.PostLinkMessageThread(r.Context(), logLink(cbState.Payload.Namespace.Name), user, channel, ts)
	}

	render.Respond(w, r, nil)
}

func (c EveController) eveCronCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the URL Params
	channel := r.URL.Query().Get("channel")

	// Get the Body
	payload := eve.NSDeploymentPlan{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		c.svc.ChatService.ErrorNotification(r.Context(), "", channel, err)
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	user := ""
	if payload.Status == eve.DeploymentPlanStatusErrors {
		user = "channel"
	}

	cbState := eveapi.CallbackState{User: user, Channel: channel, Payload: payload}
	if cbState.Payload.Status == eve.DeploymentPlanStatusPending || cbState.Payload.NothingToDeploy() {
		render.Respond(w, r, nil)
		return
	}
	ts := c.svc.ChatService.PostMessage(r.Context(), cbState.ToChatMsg(), cbState.Channel)
	if cbState.Payload.Status == eve.DeploymentPlanStatusErrors {
		c.svc.ChatService.PostLinkMessageThread(r.Context(), logLink(cbState.Payload.Namespace.Name), user, channel, ts)
	}

	render.Respond(w, r, nil)
}
