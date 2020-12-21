package main

import (
	"net/http"

	"gitlab.unanet.io/devops/eve-bot/internal/api"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	evehttp "gitlab.unanet.io/devops/go/pkg/http"
	"gitlab.unanet.io/devops/go/pkg/log"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	app, err := api.NewApi(api.InitController(&cfg), cfg)
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}

	app.Start()
}

// This is required for the HTTP Client Request/Response Logging
// Not sure why, but setting explicitly only works in the parent eve project
// when importing the mod from other repos, this needs to be handled via init process
func init() {
	http.DefaultTransport = evehttp.LoggingTransport
}
