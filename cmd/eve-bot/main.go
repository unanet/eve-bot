package main

import (
	"net/http"

	"github.com/unanet/eve-bot/internal/api"
	evehttp "github.com/unanet/go/pkg/http"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

func main() {
	app, err := api.NewApi()
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
