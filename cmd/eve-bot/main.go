package main

import (
	"gitlab.unanet.io/devops/eve-bot/internal/api"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
	"go.uber.org/zap"
)

func main() {
	api, err := mux.NewApi(api.Controllers, mux.Config{
		Port:        8080,
		MetricsPort: 3000,
		ServiceName: "eve-bot",
	})
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}
	api.Start()
}
