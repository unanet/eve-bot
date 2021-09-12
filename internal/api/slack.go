package api

import (
	"context"
	"encoding/json"
	goerror "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/unanet/eve-bot/internal/botcommander/interfaces"
	"github.com/unanet/eve-bot/internal/service"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

// SlackController for slack routes
type SlackController struct {
	svc                          *service.Provider
	exe                          interfaces.CommandExecutor
	allowedMaintenanceChannelMap map[string]interface{}
}

// NewSlackController creates a new slack controller (route handler)
func NewSlackController(svc *service.Provider, exe interfaces.CommandExecutor) *SlackController {
	return &SlackController{
		svc:                          svc,
		exe:                          exe,
		allowedMaintenanceChannelMap: extractChannelMap(svc.Cfg.SlackChannelsMaintenance),
	}
}

// Setup the routes
func (c SlackController) Setup(r chi.Router) {
	r.Post("/slack-events", c.slackEventHandler)
	r.Post("/slack-interactive", c.slackInteractiveHandler)
}

func (c SlackController) slackInteractiveHandler(w http.ResponseWriter, r *http.Request) {
	if err := handleSlackInteraction(r); err != nil {
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
	case *slack.FileSharedEvent:
		log.Logger.Info("File Uploaded", zap.Any("event", ev))
	case *slackevents.AppMentionEvent:
		c.handleSlackAppMentionEvent(r.Context(), ev)
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
	return errors.RestError{Code: status, Message: msg, OriginalError: oerr}
}

func validateSlackRequest(req *http.Request, signingSecret string) ([]byte, error) {
	verifier, err := slack.NewSecretsVerifier(req.Header, signingSecret)
	if err != nil {
		return []byte{}, botError(err, "failed new secret verifier", http.StatusUnauthorized)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return []byte{}, botError(err, "failed read req body", http.StatusBadRequest)
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

// handleSlackInteraction handles the interactive callbacks (buttons, dropdowns, etc.)
func handleSlackInteraction(req *http.Request) error {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(req.FormValue("payload")), &payload)
	if err != nil {
		return errors.RestError{Code: http.StatusBadRequest, Message: "failed to parse interactive slack message payload", OriginalError: err}
	}
	log.Logger.Info(fmt.Sprintf("Message button pressed by user %s with value %s", payload.User.Name, payload.Value))
	return nil
}

func (c SlackController) handleSlackAppMentionEvent(ctx context.Context, ev *slackevents.AppMentionEvent) {
	// Resolve the input and return an EvebotCommand object
	cmd := c.svc.CommandResolver.Resolve(ev.Text, ev.Channel, ev.User)

	chatUser, err := c.svc.ChatService.GetUser(ctx, cmd.Info().User)
	if err != nil {
		c.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, ev.ThreadTimeStamp, err)
		return
	}

	userEntry, err := c.svc.ReadUser(chatUser.FullyQualifiedName())
	if err != nil {
		if !goerror.Is(err, errors.ErrNotFound) {
			c.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, ev.ThreadTimeStamp, err)
		}
		c.svc.ChatService.PostPrivateMessage(ctx, c.svc.AuthCodeURL(chatUser.FullyQualifiedName()), cmd.Info().User)
		_ = c.svc.ChatService.PostMessageThread(ctx, "You need to login. Please Check your Private DM from `evebot` for an auth link", cmd.Info().Channel, ev.ThreadTimeStamp)
		return
	}

	// TODO: Add a "Request Access" process here when user is not authorized to perform action
	// Bonus points: Send request to configurable chat group (ex: @devops assign `username` to `role`)
	// Doubly Bonus Points: Add ability for eve to assign a user to a role
	// @evebot assign user@domain.tld to eve-deploy-prod role
	if !c.svc.IsAuthorized(cmd, userEntry) {
		_ = c.svc.ChatService.PostMessageThread(ctx, "You are not authorized to perform this action\nPlease message `@devops` with an access request if needed.", cmd.Info().Channel, ev.ThreadTimeStamp)
		return
	}

	// SlackMaintenanceEnabled is like a "feature flag"
	// set to true, and we are in Maintenance Mode
	if c.svc.Cfg.SlackMaintenanceEnabled && !userEntry.IsAdmin {
		_ = c.svc.ChatService.PostMessageThread(ctx, ":construction: Sorry, but we are currently in maintenance mode!", cmd.Info().Channel, ev.ThreadTimeStamp)
		return
	}

	// Hydrate the Acknowledgement Message and whether we should continue...
	ackMsg, cont := cmd.AckMsg()
	// Send the AckMsg and get the Timestamp back, so we can thread it later on...
	timeStamp := c.svc.ChatService.PostMessageThread(ctx, ackMsg, cmd.Info().Channel, ev.ThreadTimeStamp)
	// If the AckMessage needs to continue (no errors)...
	if cont {
		// Asynchronous CommandExecutor call
		// which maps an EveBotCommand to a CommandHandler
		go c.exe.Execute(context.TODO(), cmd, timeStamp)
	}
}

func extractChannelMap(input string) map[string]interface{} {
	chanMap := make(map[string]interface{})
	for _, c := range strings.Split(input, ",") {
		chanMap[c] = true
	}
	return chanMap
}
