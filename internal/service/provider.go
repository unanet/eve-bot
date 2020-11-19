package service

import (
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/executor"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resolver"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/config"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

// Provider provides access to the Slack Client
// and the deps required for this package
type Provider struct {
	ChatService                  chatservice.Provider
	CommandResolver              resolver.Resolver
	CommandExecutor              executor.Executor
	EveAPI                       eveapi.Client
	Cfg                          *config.Config
	allowedChannelMap            map[string]interface{}
	allowedMaintenanceChannelMap map[string]interface{}
}

// New creates a new service provider
func New(
	cfg *config.Config,
	cr resolver.Resolver,
	ea eveapi.Client,
	cs chatservice.Provider,
	ce executor.Executor,
) *Provider {

	// Elevated Slack Channels that can issue "special" commands
	// release, deploy to prod, etc.
	chanMap := make(map[string]interface{})
	for _, c := range strings.Split(cfg.SlackChannelsAuth, ",") {
		chanMap[c] = true
	}

	// Elevated Slack Channels that aren't blocked from maintenance mode
	// i.e. ops still needs to be able to test and deploy even during maintenance
	maintenanceChanMap := make(map[string]interface{})
	for _, c := range strings.Split(cfg.SlackChannelsMaintenance, ",") {
		maintenanceChanMap[c] = true
	}

	return &Provider{
		CommandResolver:              cr,
		EveAPI:                       ea,
		Cfg:                          cfg,
		ChatService:                  cs,
		CommandExecutor:              ce,
		allowedChannelMap:            chanMap,
		allowedMaintenanceChannelMap: maintenanceChanMap,
	}
}
