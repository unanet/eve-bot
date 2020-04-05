package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve-bot/internal/api/httphelper"
	"gitlab.unanet.io/devops/eve-bot/internal/api/middleware"
	"gitlab.unanet.io/devops/eve-bot/internal/api/resterror"
	"gitlab.unanet.io/devops/eve-bot/internal/servicefactory"
	"go.uber.org/zap"
)

func (a *app) registerRoutes(svcFactory *servicefactory.Container) {
	a.router.Handle("/", middleware.Handler{AppCtx: svcFactory, RouteHandler: pingHandler}).Methods("GET")
	a.router.Handle("/slack-events", middleware.Handler{AppCtx: svcFactory, RouteHandler: slackEventHandler}).Methods("POST")
}

func pingHandler(svcFactory *servicefactory.Container, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	return httphelper.AppResponse(http.StatusOK, "pong")
}

func slackEventHandler(svcFactory *servicefactory.Container, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {

	cleanErr := func(oerr error, msg string, status int) (error, string) {
		return &resterror.RestError{
			Code:          http.StatusUnauthorized,
			Message:       msg,
			OriginalError: oerr,
		}, msg
	}

	verifier, err := slack.NewSecretsVerifier(req.Header, svcFactory.Config.SlackSecrets.SigningSecret)
	if err != nil {
		return httphelper.AppErr(cleanErr(err, "failed new secret verifier", http.StatusUnauthorized))
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return httphelper.AppErr(cleanErr(err, "failed readAll req body", http.StatusBadRequest))
	}

	_, err = verifier.Write(body)
	if err != nil {
		return httphelper.AppErr(cleanErr(err, "failed verifier write", http.StatusUnauthorized))
	}

	err = verifier.Ensure()
	if err != nil {
		return httphelper.AppErr(cleanErr(err, "failed verifier ensure", http.StatusUnauthorized))
	}

	slackAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: svcFactory.Config.SlackSecrets.VerificationToken}))
	if err != nil {
		return httphelper.AppErr(err, "failed parse slack event")
	}

	svcFactory.Logger.For(req.Context()).Debug("slack event", zap.String("event", slackAPIEvent.Type))

	if slackAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return httphelper.AppErr(err, "slack reg events url verification")
		}
		return httphelper.AppResponse(http.StatusOK, []byte(r.Challenge))
	}

	if slackAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// TODO: This is where we are going to handle the incoming text from the User
			// user := ev.User
			// chatMsg := ev.Text

			svcFactory.Logger.For(req.Context()).Debug("slack event Data >>>>>>>>>>>>>>>",
				zap.String("type", slackAPIEvent.Type),
				zap.Any("data", slackAPIEvent.Data),
				zap.String("ev.User", ev.User),
				zap.String("ev.Text", ev.Text),
				zap.Any("innerEvent.Data", innerEvent.Data),
				zap.Any("innerEvent.Type", innerEvent.Type),
			)

			msgFields := strings.Fields(ev.Text)

			botIDField := msgFields[0]
			commandFields := msgFields[1:]

			// attachment := slack.Attachment{}

			slack.MsgOptionAttachments(slack.Attachment{})

			svcFactory.SlackClient.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Yes, <@%s>, this is %s! You want me to run: %s", ev.User, botIDField, commandFields), false))
			return httphelper.AppResponse(http.StatusOK, "")
		}
	}

	return httphelper.AppResponse(http.StatusGone, "unknown slack event")

}
