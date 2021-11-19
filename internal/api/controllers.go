package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-chi/chi"
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/commands/handlers"
	"github.com/unanet/eve-bot/internal/botcommander/executor"
	"github.com/unanet/eve-bot/internal/botcommander/resolver"
	chat "github.com/unanet/eve-bot/internal/chatservice"
	"github.com/unanet/eve-bot/internal/config"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/eve-bot/internal/service"
	"github.com/unanet/go/pkg/identity"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

type Controller interface {
	Setup(chi.Router)
}

// initController initializes the controller (handlers)
func initController(cfg *config.Config) []Controller {
	eveAPI := eveapi.New(cfg.EveAPIConfig)
	chatSvc := chat.New(chat.Slack, cfg)

	awsSession, err := session.NewSession(&aws.Config{Region: aws.String(cfg.AWSRegion)})
	if err != nil {
		log.Logger.Panic("Unable to Initialize the AWS Session", zap.Error(err))
	}

	idSvc, err := identity.NewValidator(cfg.Identity)
	if err != nil {
		log.Logger.Panic("Unable to Initialize the Identity Service Provider", zap.Error(err))
	}

	svc := service.New(cfg,
		service.ChatProviderParam(chatSvc),
		service.DynamoParam(dynamodb.New(awsSession)),
		service.EveAPIParam(eveAPI),
		service.ResolverParam(resolver.New(commands.NewFactory())),
		service.OpenIDConnectParam(cfg, idSvc),
	)

	exe := executor.New(svc, handlers.NewFactory())

	return []Controller{
		NewPingController(),
		NewSlackController(svc, exe),
		NewEveController(svc),
		NewAuthController(svc),
	}
}
