package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/api/controller"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/eveerrs"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type Controller struct {
	controller.Base
	slackProvider *islack.Provider
}

func New(slackProvider *islack.Provider) *Controller {
	return &Controller{
		slackProvider: slackProvider,
	}
}

func (c Controller) Setup(r chi.Router) {
	r.Post("/slack-events", c.slackEventHandler)
}

func (c Controller) slackEventHandler(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Slack Event Handler")
	// err := c.slackProvider.HandleEvent(r)

	body, err := verifyRequestSig(r, &config.Values().SlackSigningSecret)

	if err != nil {
		log.Logger.Error("Slack Event Handler Error verify request sig", zap.Error(err))
		// render.Respond(w, r, &eveerrs.RestError{
		// 	Code:          400,
		// 	Message:       "failed slack verification",
		// 	OriginalError: nil,
		// }) OR just send an error back if this is temp
		render.Respond(w, r, fmt.Errorf("failed slack verification"))
		return
	}

	slackAPIEvent, err := slackevents.ParseEvent(
		body,
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: config.Values().SlackVerificationToken}))
	if err != nil {
		log.Logger.Error("Slack Event Handler Error ParseEvent", zap.Error(err))
		render.Respond(w, r, fmt.Errorf("failed parse slack event"))
		return
	}

	log.Logger.Info("slack event parsed", zap.String("event", slackAPIEvent.Type))

	if slackAPIEvent.Type == slackevents.URLVerification {
		var cr *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &cr)
		if err != nil {
			render.Respond(w, r, fmt.Errorf("slack reg url verification"))
			return
		}
		render.Respond(w, r, []byte(cr.Challenge))
		return
	}

	if slackAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// TODO: This is where we are going to handle the incoming text from the User
			// user := ev.User
			// chatMsg := ev.Text

			// svcFactory.Logger.For(req.Context()).Debug("slack event Data >>>>>>>>>>>>>>>",
			// 	zap.String("type", slackAPIEvent.Type),
			// 	zap.Any("data", slackAPIEvent.Data),
			// 	zap.String("ev.User", ev.User),
			// 	zap.String("ev.Text", ev.Text),
			// 	zap.Any("innerEvent.Data", innerEvent.Data),
			// 	zap.Any("innerEvent.Type", innerEvent.Type),
			// )

			msgFields := strings.Fields(ev.Text)

			botIDField := msgFields[0]
			commandFields := msgFields[1:]

			fmt.Printf("%s, %s", botIDField, commandFields)

			// attachment := slack.Attachment{}

			slack.MsgOptionAttachments(slack.Attachment{})

			c.slackProvider.Client.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Yes, <@%s>, this is %s! You want me to run: %s", ev.User, botIDField, commandFields), false))
			render.Respond(w, r, "")
			return
		}
	}

	render.Status(r, http.StatusGone)
	render.Respond(w, r, "unknown slack event")

	if err != nil {
		log.Logger.Error("Slack Event Handler Error", zap.Error(err))
		render.Respond(w, r, &eveerrs.RestError{
			Code:          400,
			Message:       "Bad Request",
			OriginalError: err,
		})
	}

	render.Respond(w, r, nil)
}

// Private/Helper functions
func verifyRequestSig(req *http.Request, signingSecret *string) ([]byte, error) {
	cleanErr := func(oerr error, msg string, status int) error {
		log.Logger.Error("verifyRequestSig error", zap.Error(oerr))
		return &eveerrs.RestError{
			Code:          http.StatusUnauthorized,
			Message:       msg,
			OriginalError: oerr,
		}
	}

	log.Logger.Info("Request Header TS:", zap.String("val", req.Header.Get("X-Slack-Request-Timestamp")))
	log.Logger.Info("Request Header Sig:", zap.String("val", req.Header.Get("X-Slack-Signature")))

	verifier, err := slack.NewSecretsVerifier(req.Header, *signingSecret)
	if err != nil {
		return []byte{}, cleanErr(err, "failed new secret verifier", http.StatusUnauthorized)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return []byte{}, cleanErr(err, "failed readAll req body", http.StatusBadRequest)
	}

	_, err = verifier.Write(body)
	if err != nil {
		return []byte{}, cleanErr(err, "failed verifier write", http.StatusUnauthorized)
	}

	err = verifier.Ensure()
	if err != nil {
		return []byte{}, cleanErr(err, "failed verifier ensure", http.StatusUnauthorized)
	}

	return body, nil
}
