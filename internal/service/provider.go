package service

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/unanet/go/pkg/identity"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"

	"github.com/unanet/eve-bot/internal/config"
)

// Provider provides access to the Common Deps/Service required for this project
type Provider struct {
	ChatService                  interfaces.ChatProvider
	CommandResolver              interfaces.CommandResolver
	CommandExecutor              interfaces.CommandExecutor
	EveAPI                       interfaces.EveAPI
	Cfg                          *config.Config
	allowedMaintenanceChannelMap map[string]interface{}
	mutex                        sync.Mutex
	oidc                         *identity.Service
	userDB                       *dynamodb.DynamoDB
	userCache                    map[string]UserEntry
}

func OpenIDConnectParam(id *identity.Service) Option {
	return func(svc *Provider) {
		svc.oidc = id
	}
}

func DynamoParam(db *dynamodb.DynamoDB) Option {
	return func(svc *Provider) {
		svc.userDB = db
	}
}

func ResolverParam(r interfaces.CommandResolver) Option {
	return func(svc *Provider) {
		svc.CommandResolver = r
	}
}

func EveAPIParam(e interfaces.EveAPI) Option {
	return func(svc *Provider) {
		svc.EveAPI = e
	}
}

func ChatProviderParam(c interfaces.ChatProvider) Option {
	return func(svc *Provider) {
		svc.ChatService = c
	}
}

func ExecutorParam(c interfaces.CommandExecutor) Option {
	return func(svc *Provider) {
		svc.CommandExecutor = c
	}
}

type Option func(*Provider)

func New(cfg *config.Config, opts ...Option) *Provider {
	svc := &Provider{
		Cfg: cfg,
		// Elevated Slack Channels that aren't blocked from maintenance mode
		// i.e. ops still needs to be able to test and deploy even during maintenance
		allowedMaintenanceChannelMap: extractChannelMap(cfg.SlackChannelsMaintenance),
		mutex:                        sync.Mutex{},
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}
