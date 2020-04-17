package ping

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve-bot/internal/api/controller"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

type Controller struct {
	controller.Base
}

func New() *Controller {
	return &Controller{}
}

func (c Controller) Setup(r chi.Router) {
	r.Get("/ping", c.ping)
}

func (c Controller) ping(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, render.M{
		"message": "pong",
		"version": mux.Version,
	})
}
