package main

import (
	"gitlab.unanet.io/devops/eve-bot/internal/api"
	"gitlab.unanet.io/devops/eve-bot/internal/config"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
	"go.uber.org/zap"
)

// adding a comment to test deploys
func main() {

	//queue.StartDispatcher(config.Values().QueueWorkers)

	app, err := mux.NewApi(api.Controllers, config.Values().MuxConfig)
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}
	//app.Start(func() {
	//
	//})

	app.Start()

}
