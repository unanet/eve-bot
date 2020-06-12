package evebotservice

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/executor"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resolver"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/errors"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// Provider provides access to the Slack Client
// and the deps required for this package
type Provider struct {
	ChatService     chatservice.Provider
	CommandResolver resolver.Resolver
	CommandExecutor executor.Executor
	EveAPI          eveapi.Client
	Cfg             *config.Config
}

// NewProvider creates a new provider
func New(
	cfg *config.Config,
	cr resolver.Resolver,
	ea eveapi.Client,
	cs chatservice.Provider,
	ce executor.Executor,
) *Provider {
	return &Provider{
		CommandResolver: cr,
		EveAPI:          ea,
		Cfg:             cfg,
		ChatService:     cs,
		CommandExecutor: ce,
	}
}

func botError(oerr error, msg string, status int) error {
	log.Logger.Debug("EveBot Error", zap.Error(oerr))
	return &errors.RestError{Code: status, Message: msg, OriginalError: oerr}
}
