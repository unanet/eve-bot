package servicefactory

import (
	"github.com/slack-go/slack"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/evelogger"
	"gitlab.unanet.io/devops/eve-bot/internal/metrics"
	"gitlab.unanet.io/devops/eve-bot/internal/version"
	"go.uber.org/zap"
)

// Container contains the shared services
type Container struct {
	VersionInfo version.Info
	Config      *config.Config
	Metrics     *metrics.Provider
	Logger      evelogger.Container
	SlackClient *slack.Client
}

var zlogger *zap.Logger

// Initialize is the main bootstrap/initialization process
// if this fails, we should fail hard (panic/fatal etc.)
// this should only be called once, on app startup
// NEW is glue and this is where that glue is applied
// create your deps here and add them to the Container
func Initialize(vInfo version.Info) *Container {
	cfg := config.Read()

	zlogger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger")
	}

	logger := evelogger.NewLogContainer(cfg, zlogger.With(zap.String("package", "servicefactory")))

	return &Container{
		VersionInfo: vInfo,
		Config:      cfg,
		Metrics:     metrics.New(),
		Logger:      logger,
		SlackClient: slack.New(cfg.SlackSecrets.BotOAuthToken),
	}
}
