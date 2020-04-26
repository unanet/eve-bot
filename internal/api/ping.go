package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"gitlab.unanet.io/devops/eve/pkg/mux"
)

type PingController struct {
}

func NewPingController() *PingController {
	return &PingController{}
}

func (c PingController) Setup(r chi.Router) {
	r.Get("/ping", c.ping)
}

func (c PingController) ping(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, render.M{
		"message": "pong",
		"version": mux.Version,
	})
}
