package service

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/unanet/eve-bot/internal/manager"
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"

	"github.com/unanet/eve-bot/internal/config"
)

// Provider provides access to the Slack Client
// and the deps required for this package
type Provider struct {
	ChatService                  interfaces.ChatProvider
	CommandResolver              interfaces.CommandResolver
	CommandExecutor              interfaces.CommandExecutor
	EveAPI                       interfaces.EveAPI
	UserDB                       *dynamodb.DynamoDB
	Cfg                          *config.Config
	MgrSvc                       *manager.Service
	allowedChannelMap            map[string]interface{}
	allowedMaintenanceChannelMap map[string]interface{}
}

func extractChannelMap(input string) map[string]interface{} {
	chanMap := make(map[string]interface{})
	for _, c := range strings.Split(input, ",") {
		chanMap[c] = true
	}
	return chanMap
}

// New creates a new service provider
func New(
	cfg *config.Config,
	cr interfaces.CommandResolver,
	ea interfaces.EveAPI,
	cs interfaces.ChatProvider,
	ce interfaces.CommandExecutor,
	svc *dynamodb.DynamoDB,
	mgr *manager.Service) *Provider {

	return &Provider{
		MgrSvc:          mgr,
		UserDB:          svc,
		CommandResolver: cr,
		EveAPI:          ea,
		Cfg:             cfg,
		ChatService:     cs,
		CommandExecutor: ce,
		// Elevated Slack Channels that can issue "special" commands
		// release, deploy to prod, etc.
		allowedChannelMap: extractChannelMap(cfg.SlackChannelsAuth),
		// Elevated Slack Channels that aren't blocked from maintenance mode
		// i.e. ops still needs to be able to test and deploy even during maintenance
		allowedMaintenanceChannelMap: extractChannelMap(cfg.SlackChannelsMaintenance),
	}
}
