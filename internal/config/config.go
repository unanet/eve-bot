package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/unanet/eve-bot/internal/chatservice/slackservice"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/go/pkg/identity"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

var (
	cfg   *Config
	mutex = sync.Mutex{}
)

type (
	// LogConfig is the logger config (log level, output...)
	LogConfig = log.Config
	// SlackConfig is the slack config (secret, tokens...)
	SlackConfig = slackservice.Config
	// EveAPIConfig is the config for the Eve API
	EveAPIConfig = eveapi.Config
	// IdentityConfig is the OIDC (KeyCloak) Config data
	IdentityConfig = identity.Config
)

// Config is the top level application config
type Config struct {
	LogConfig
	SlackConfig
	EveAPIConfig
	Identity                IdentityConfig
	Port                    int    `split_words:"true" default:"8080"`
	MetricsPort             int    `split_words:"true" default:"3001"`
	ServiceName             string `split_words:"true" default:"eve"`
	ReadOnly                bool   `split_words:"true" default:"false"`
	AWSRegion               string `split_words:"true" required:"true"`
	LoggingDashboardBaseURL string `split_words:"true" required:"true"`
}

// Load loads the config reading it from the environment
func Load() Config {
	mutex.Lock()
	defer mutex.Unlock()
	if cfg != nil {
		return *cfg
	}
	c := Config{}
	err := envconfig.Process("EVEBOT", &c)
	if err != nil {
		log.Logger.Panic("Unable to Load Config", zap.Error(err))
	}
	cfg = &c
	return *cfg
}
