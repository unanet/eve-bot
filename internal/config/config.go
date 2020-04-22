package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	islack "gitlab.unanet.io/devops/eve-bot/internal/slack"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"gitlab.unanet.io/devops/eve/pkg/mux"
	"go.uber.org/zap"
)

var (
	values *Config
	mutex  = sync.Mutex{}
)

type (
	// LogConfig is the logger config (log level, output...)
	LogConfig = log.Config
	// SlackConfig is the slack config (secret, tokens..)
	SlackConfig = islack.Config
	// MuxConfig is the multiplexer (router) config (ports, timeouts)
	MuxConfig = mux.Config
)

// Config is the top level application config
type Config struct {
	LogConfig
	SlackConfig
	MuxConfig
}

// Values returns the environmental config values (prefix: EVEBOT_)
func Values() *Config {
	mutex.Lock()
	defer mutex.Unlock()
	if values != nil {
		return values
	}
	c := Config{}
	err := envconfig.Process("EVEBOT", &c)
	if err != nil {
		log.Logger.Panic("Unable to Load Config", zap.Error(err))
	}
	values = &c
	return values
}
