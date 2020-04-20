package main

import (
	"gitlab.unanet.io/devops/eve-bot/internal/api"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
	"go.uber.org/zap"
)

func main() {
	api, err := mux.NewApi(api.Controllers, mux.Config{
		Port:        3000,
		MetricsPort: 8080,
		ServiceName: "eve-bot",
	})
	if err != nil {
		log.Logger.Panic("Failed to Create Api App", zap.Error(err))
	}
	api.Start()
}
