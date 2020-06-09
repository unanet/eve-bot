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
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/evebotservice"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// Controller for slack routes
type SlackController struct {
	svc *evebotservice.Provider
}

// New creates a new slack controller (route handler)
func NewSlackController(svc *evebotservice.Provider) *SlackController {
	return &SlackController{svc: svc}
}

// Setup the routes
func (c SlackController) Setup(r chi.Router) {
	r.Post("/slack-events", c.slackEventHandler)
	r.Post("/slack-interactive", c.slackInteractiveHandler)
	r.Post("/eve-callback", c.eveCallbackHandler)
	r.Post("/eve-cron-callback", c.eveCronCallbackHandler)
}

func logLink(ns string) string {
	return "https://grafana.unanet.io/explore?left=%5B%22now-1h%22,%22now%22,%22Loki%22,%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D~%5C%22" + ns + ".*%5C%22%7D%22,%22key%22:%22Q-1591690807784-0.2580443768677896-0%22,%22hide%22:false%7D,%7B%22mode%22:%22Logs%22%7D,%7B%22ui%22:%5Btrue,true,true,null%5D%7D%5D"
}

func (c SlackController) eveCallbackHandler(w http.ResponseWriter, r *http.Request) {
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
	log.Logger.Debug("eve callback notification", zap.Any("cb_state", cbState))
	c.svc.ChatService.PostMessageThread(r.Context(), cbState.ToChatMsg(), cbState.Channel, cbState.TS)

	if cbState.Payload.Status == eve.DeploymentPlanStatusErrors {
		c.svc.ChatService.PostLinkMessageThread(r.Context(), logLink(cbState.Payload.Namespace.Name), user, channel, ts)
	}

	render.Respond(w, r, nil)
}

func (c SlackController) eveCronCallbackHandler(w http.ResponseWriter, r *http.Request) {
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
	log.Logger.Debug("eve cron callback notification", zap.Any("cb_state", cbState))
	if cbState.Payload.Status == eve.DeploymentPlanStatusPending {
		render.Respond(w, r, nil)
		return
	}
	c.svc.ChatService.PostMessage(r.Context(), cbState.ToChatMsg(), cbState.Channel)
	render.Respond(w, r, nil)
	return

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

	// We are only handling/listening to the CallbackEvent
	if slackAPIEvent.Type != slackevents.CallbackEvent {
		render.Respond(w, r, errors.Wrap(fmt.Errorf("unknown slack API event: %s", slackAPIEvent.Type)))
		return
	}
	innerEvent := slackAPIEvent.InnerEvent
	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		if err := c.svc.HandleSlackAppMentionEvent(r.Context(), ev); err != nil {
			render.Respond(w, r, errors.Wrap(err))
			return
		}
	default:
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
