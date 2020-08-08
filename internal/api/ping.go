package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/pkg/mux"
)

// PingController is the Controller/Handler for ping routes
type PingController struct {
}

// NewPingController creates a new ping controller
func NewPingController() *PingController {
	return &PingController{}
}

// Setup satisfies the EveController interface for setting up the
func (c PingController) Setup(r chi.Router) {
	r.Get("/ping", c.ping)
}

func (c PingController) ping(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, render.M{
		"message": "pong",
		"version": mux.Version,
	})
}
