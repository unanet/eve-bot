package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice/slackservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
	"go.uber.org/zap"
)

var (
	cfg   *Config
	mutex = sync.Mutex{}
)

type (
	// LogConfig is the logger config (log level, output...)
	LogConfig = log.Config
	// SlackConfig is the slack config (secret, tokens..)
	SlackConfig = slackservice.Config
	// MuxConfig is the multiplexer (router) config (ports, timeouts)
	MuxConfig = mux.Config
	// EveAPIConfig is the config for the Eve API
	EveAPIConfig = eveapi.Config
)

// Config is the top level application config
type Config struct {
	LogConfig
	SlackConfig
	MuxConfig
	EveAPIConfig
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
