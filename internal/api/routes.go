package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve-bot/internal/api/httphelper"
	"gitlab.unanet.io/devops/eve-bot/internal/api/middleware"
	"gitlab.unanet.io/devops/eve-bot/internal/servicefactory"
	"go.uber.org/zap"
)

func (a *app) registerRoutes(svcFactory *servicefactory.Container) {
	a.router.Handle("/", middleware.Handler{AppCtx: svcFactory, RouteHandler: pingHandler}).Methods("GET")
	a.router.Handle("/ping", middleware.Handler{AppCtx: svcFactory, RouteHandler: pingHandler}).Methods("GET")
	a.router.Handle("/slack-event", middleware.Handler{AppCtx: svcFactory, RouteHandler: slackEventHandler}).Methods("GET")
}

func pingHandler(svcFactory *servicefactory.Container, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {
	return httphelper.AppResponse(http.StatusOK, "pong")
}

func slackEventHandler(svcFactory *servicefactory.Container, res http.ResponseWriter, req *http.Request) (int, interface{}, error) {

	verifier, err := slack.NewSecretsVerifier(req.Header, svcFactory.Config.SlackSecrets.SigningSecret)
	if err != nil {
		return httphelper.AppErr(err, "failed new secret verifier")
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return httphelper.AppErr(err, "failed readAll req body")
	}

	_, err = verifier.Write(body)
	if err != nil {
		return httphelper.AppErr(err, "failed verifier write")
	}

	err = verifier.Ensure()
	if err != nil {
		return httphelper.AppErr(err, "failed verifier ensure")
	}

	slackAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: svcFactory.Config.SlackSecrets.VerificationToken}))
	if err != nil {
		return httphelper.AppErr(err, "failed parse slack event")
	}

	if slackAPIEvent.Type == slackevents.URLVerification {
		svcFactory.Logger.For(req.Context()).Info("slack event", zap.String("event", slackevents.URLVerification))
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return httphelper.AppErr(err, "slack reg events url verification")
		}
		return httphelper.AppResponse(http.StatusOK, []byte(r.Challenge))
	}

	if slackAPIEvent.Type == slackevents.CallbackEvent {
		svcFactory.Logger.For(req.Context()).Info("slack event", zap.String("event", slackevents.CallbackEvent))
		innerEvent := slackAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			svcFactory.SlackClient.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			return httphelper.AppResponse(http.StatusOK, "")
		}
	}

	svcFactory.Logger.For(req.Context()).Info("slack event", zap.String("event", "unknown"))
	return httphelper.AppResponse(http.StatusOK, "unknown slack event")

}
