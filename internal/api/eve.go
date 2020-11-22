package api

import (
	"encoding/json"
	"net/http"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.unanet.io/devops/eve-bot/internal/service"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
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
	return "https://grafana.unanet.io/explore?left=%5B%22now-1h%22,%22now%22,%22Loki%22,%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D~%5C%22" + ns + ".*%5C%22%7D%22,%22key%22:%22Q-1591690807784-0.2580443768677896-0%22,%22hide%22:false%7D,%7B%22mode%22:%22Logs%22%7D,%7B%22ui%22:%5Btrue,true,true,null%5D%7D%5D"
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
