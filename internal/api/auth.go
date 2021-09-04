package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/unanet/eve-bot/internal/service"
	"github.com/unanet/go/pkg/errors"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
	"net/http"
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
}

func (c AuthController) successfulSignIn(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("<!doctype html>\n\n<html lang=\"en\">\n<head>\n <script language=\"javascript\" type=\"text/javascript\">\nfunction windowClose() {\nwindow.open('','_parent','');\nwindow.close();\n}\n</script> <meta charset=\"utf-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>Successful Auth</title>\n</head>\n<body>\n  \t<p> You have successfully Signed In. You may close this window</p>\n</body>\n</html>"))
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
