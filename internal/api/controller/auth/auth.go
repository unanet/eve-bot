package auth

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"gitlab.unanet.io/devops/eve-bot/internal/api/controller"
)

// New returns a new Auth Controller (ctrl)
func New() *Controller {
	return &Controller{}
}

// Controller is the Auth Controller
type Controller struct {
	controller.Base
}

// Setup implements the Eve Controller Interface
func (c *Controller) Setup(r chi.Router) {
	r.Route("/authorization-code/callback", func(r chi.Router) {
		r.Post("/", c.authCodeCBHandler)
		r.Get("/", c.authCodeCBHandler)
	})
	r.Get("/logout", c.logoutHandler)
}

func (c *Controller) logoutHandler(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, "logout")
}

func (c *Controller) authCodeCBHandler(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, "auth code")
}
