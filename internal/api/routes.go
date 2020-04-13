package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/slack-go/slack"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"github.com/slack-go/slack/slackevents"
	"gitlab.unanet.io/devops/eve/pkg/eveerrs"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/config"
)

var Controllers = []mux.EveController{
	New(islack.NewProvider(config.Values)),
}

type MyController struct {
	slackProvider islack.Provider
}

func (c MyController) Setup(r chi.Router) {
	r.Get("/", c.pingHandler)
	r.Post("/slack-events", c.slackEventHandler)
	r.Route("/authorization-code/callback", func(r chi.Router) {
		r.Post("/", c.authCodeCBHandler)
		r.Get("/", c.authCodeCBHandler)
	})
	// r.Get("/login", c.loginHandler)
	r.Get("/logout", c.logoutHandler)
}

func New(slackProvider islack.Provider) *MyController {
	return &MyController{
		slackProvider: slackProvider,
	}
}




var state = "ApplicationState"
var nonce = "NonceNotSetYet"

// Handler Functions

// func loginHandler(w http.ResponseWriter, req *http.Request) {
// 	nonce, _ = httphelper.GenerateNonce()
// 	//var redirectPath string

// 	q := req.URL.Query()
// 	q.Add("client_id", svcFactory.Config.OktaSecrets.ClientID)
// 	q.Add("response_type", "code")
// 	q.Add("response_mode", "query")
// 	q.Add("scope", "openid profile email")
// 	q.Add("redirect_uri", "http://localhost:3000/authorization-code/callback")
// 	q.Add("state", state)
// 	q.Add("nonce", nonce)

// 	//redirectPath = svcFactory.Config.OktaSecrets.IssuerURL + "/v1/authorize?" + q.Encode()

// 	//svcFactory.Logger.Bg().Fatal("HELLLO", zap.String("redir_url", redirectPath))
// 	// svcFactory.Logger.For(req.Context()).Fatal("HELLLO", zap.String("redir_url", redirectPath))

// 	http.Redirect(res, req, "https://google.com", http.StatusMovedPermanently)
// 	return httphelper.AppResponse(http.StatusOK, "login")
// }

func (c MyController) logoutHandler(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, "logout")
}

func (c MyController) authCodeCBHandler(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, "auth code")
}

func (c MyController) pingHandler(w http.ResponseWriter, r *http.Request){
	render.Respond(w, r, "pong")
}

func (c MyController) slackEventHandler(w http.ResponseWriter, r *http.Request) {

	body, err := verifyRequestSig(r, &config.Values.SlackSecrets.SigningSecret)

	if err != nil {
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
		slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: config.Values.SlackSecrets.VerificationToken}))
	if err != nil {
		render.Respond(w, r, fmt.Errorf("failed parse slack event"))
		return
	}

	log.Logger.Debug("slack event", zap.String("event", slackAPIEvent.Type))

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

			//TODO: FROM CASEY on't know what you're doing HERE
			// c.slackProvider.PostMessage(ev.Channel, slack.MsgOptionText(fmt.Sprintf("Yes, <@%s>, this is %s! You want me to run: %s", ev.User, botIDField, commandFields), false))
			render.Respond(w, r, "")
			return
		}
	}

	render.Status(r, http.StatusGone)
	render.Respond(w, r, "unknown slack event")
}

// Private/Helper functions
func verifyRequestSig(req *http.Request, signingSecret *string) ([]byte, error) {
	cleanErr := func(oerr error, msg string, status int) error {
		return &eveerrs.RestError{
			Code:          http.StatusUnauthorized,
			Message:       msg,
			OriginalError: oerr,
		}
	}

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
