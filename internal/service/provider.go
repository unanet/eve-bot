package service

import (
	"context"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/coreos/go-oidc"
	"github.com/unanet/go/pkg/identity"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"

	"github.com/unanet/eve-bot/internal/config"
)

// Provider provides access to the Common Deps/Services required for this project
type Provider struct {
	ChatService     interfaces.ChatProvider
	CommandResolver interfaces.CommandResolver
	EveAPI          interfaces.EveAPI
	Cfg             *config.Config
	oidc            *identity.Validator
	userDB          *dynamodb.DynamoDB
	oauth           struct {
		config   oauth2.Config
		verifier *oidc.IDTokenVerifier
	}
}

func OpenIDConnectParam(cfg *config.Config, id *identity.Validator) Option {

	oauthProvider, err := oidc.NewProvider(context.Background(), cfg.Identity.ConnectionURL)
	if err != nil {
		log.Logger.Panic("Unable to Initialize the OAuth Provider", zap.Error(err))
	}

	return func(svc *Provider) {
		svc.oidc = id

		svc.oauth.verifier = oauthProvider.Verifier(&oidc.Config{
			ClientID: cfg.Identity.ClientID,
		})

		svc.oauth.config = oauth2.Config {
			ClientID:     cfg.Identity.ClientID,
			ClientSecret: cfg.Oidc.ClientSecret,
			RedirectURL:  cfg.Oidc.RedirectURL,
			Endpoint:     oauthProvider.Endpoint(),
			// "openid" is a required scope for OpenID Connect flows.
			Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
		}
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

type Option func(*Provider)

func New(cfg *config.Config, opts ...Option) *Provider {
	svc := &Provider{
		Cfg: cfg,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}
