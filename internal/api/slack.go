package api

import (
	"net/http"

	"gitlab.unanet.io/devops/eve-bot/internal/botmetrics"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/eveerrs"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

// Controller for slack routes
type SlackController struct {
	Base
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
	err := c.slackProvider.HandleEveCallback(r)

	if err != nil {
		// This is Bad and we should get paged
		// if this hits, there is no way to notify the user in slack
		log.Logger.Error("eve callback handler error", zap.Error(err))
		botmetrics.StatBotErrCount.WithLabelValues("eve-callback").Inc()
		render.Respond(w, r, &eveerrs.RestError{
			Code:          500,
			Message:       "unknown eve callback error",
			OriginalError: err,
		})
	}

	// Just returning an empty response here...
	render.Respond(w, r, nil)

}

func (c SlackController) slackInteractiveHandler(w http.ResponseWriter, r *http.Request) {
	err := c.slackProvider.HandleInteraction(r)

	if err != nil {
		// This is Bad and we should get paged
		// if this hits, there is no way to notify the user in slack
		log.Logger.Error("slack interaction handler error", zap.Error(err))
		botmetrics.StatBotErrCount.WithLabelValues("slack-interactive").Inc()
		render.Respond(w, r, &eveerrs.RestError{
			Code:          500,
			Message:       "unknown slack interaction error",
			OriginalError: err,
		})
	}

	// Just returning an empty response here...
	render.Respond(w, r, nil)
}

func (c SlackController) slackEventHandler(w http.ResponseWriter, r *http.Request) {
	err := c.slackProvider.HandleEvent(r)

	if err != nil {
		// This is a Bad scenario and we should get paged
		// if this hits, there is no way to notify the user in slack
		log.Logger.Error("slack event handler error", zap.Error(err))
		botmetrics.StatBotErrCount.WithLabelValues("slack-event").Inc()
		render.Respond(w, r, &eveerrs.RestError{
			Code:          500,
			Message:       "unknown slack event error",
			OriginalError: err,
		})
	}

	// Just returning an empty response here...
	render.Respond(w, r, nil)
}
