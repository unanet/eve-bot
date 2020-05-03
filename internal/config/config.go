package config

import (
	"sync"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

	"github.com/kelseyhightower/envconfig"
	"gitlab.unanet.io/devops/eve-bot/internal/queue"
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
	// QueueConfig is the config for queue for Async Eve-API Calls
	QueueConfig = queue.Config
	// EveAPIConfig is the config for the Eve API
	EveAPIConfig = eveapi.Config
)

// Config is the top level application config
type Config struct {
	LogConfig
	SlackConfig
	MuxConfig
	QueueConfig
	EveAPIConfig
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
