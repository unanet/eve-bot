package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/unanet/eve-bot/internal/service"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
	"github.com/unanet/go/pkg/middleware"
	"go.uber.org/zap"
)

// AuthController is the Controller/Handler for ping routes
type AuthController struct {
	svc *service.Provider
}

// NewAuthController creates a new OIDC controller
func NewAuthController(svc *service.Provider) *AuthController {
	return &AuthController{
		svc: svc,
	}
}

// Setup satisfies the EveController interface for setting up the
func (c AuthController) Setup(r chi.Router) {
	r.Get("/oidc/callback", c.callback)
	r.Get("/signed-in", c.successfulSignIn)
	r.Get("/auth", c.auth)
}

func (c AuthController) successfulSignIn(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("<!doctype html>\n\n<html lang=\"en\">\n<head>\n <script language=\"javascript\" type=\"text/javascript\">\nfunction windowClose() {\nwindow.open('','_parent','');\nwindow.close();\n}\n</script> <meta charset=\"utf-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>Successful Auth</title>\n</head>\n<body>\n  \t<p> You have successfully Signed In. You may close this window</p>\n</body>\n</html>"))
}

func (c AuthController) auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	unknownToken := jwtauth.TokenFromHeader(r)

	if len(unknownToken) == 0 {
		middleware.Log(ctx).Debug("unknown token")
		http.Redirect(w, r, c.svc.AuthCodeURL("tsampson"), http.StatusFound)
		return
	}

	verifiedToken, err := c.svc.Verify(ctx, unknownToken)
	if err != nil {
		middleware.Log(ctx).Debug("invalid token")
		http.Redirect(w, r, c.svc.AuthCodeURL("tsampson"), http.StatusFound)
		return
	}

	var idTokenClaims = new(json.RawMessage)
	if err := verifiedToken.Claims(&idTokenClaims); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, TokenResponse{
		AccessToken: unknownToken,
		Expiry:      verifiedToken.Expiry,
		Claims:      idTokenClaims,
	})
}

type TokenResponse struct {
	AccessToken string
	Expiry      time.Time
	Claims      *json.RawMessage
}

func (c AuthController) callback(w http.ResponseWriter, r *http.Request) {
	incomingState := r.URL.Query().Get("state")
	log.Logger.Info("incoming oidc callback state", zap.Any("state", incomingState))
	ctx := r.Context()

	err := c.svc.SaveUserAuth(ctx,
		r.URL.Query().Get("state"),
		r.URL.Query().Get("code"))
	if err != nil {
		render.Respond(w, r, errors.Wrap(err))
		return
	}

	// Just redirecting to a different page to prevent id refresh (which throws an error)
	http.Redirect(w, r, "/signed-in", http.StatusFound)
	return
}
