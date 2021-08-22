package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	//"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-chi/chi"
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/commands/handlers"
	"github.com/unanet/eve-bot/internal/botcommander/executor"
	"github.com/unanet/eve-bot/internal/botcommander/resolver"
	chat "github.com/unanet/eve-bot/internal/chatservice"
	"github.com/unanet/eve-bot/internal/config"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/eve-bot/internal/manager"
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

	cmdResolver := resolver.New(commands.NewFactory())
	eveAPI := eveapi.New(cfg.EveAPIConfig)
	chatSvc := chat.New(chat.Slack, cfg)
	cmdExecutor := executor.New(eveAPI, chatSvc, handlers.NewFactory())

	awsSession, err := session.NewSession(&aws.Config{Region: aws.String(cfg.AWSRegion)})
	if err != nil {
		log.Logger.Panic("Unable to Initialize the AWS Session", zap.Error(err))
	}

	// Create the Service Deps here
	idSvc, err := identity.NewService(cfg.Identity)
	if err != nil {
		log.Logger.Panic("Unable to Initialize the Identity Service Provider", zap.Error(err))
	}

	// Create the Service Managers here
	mgr := manager.NewService(cfg, manager.OpenIDConnectOpt(idSvc))

	// Create DynamoDB client
	dynamoDBSvc := dynamodb.New(awsSession)

	svc := service.New(
		cfg,
		cmdResolver,
		eveAPI,
		chatSvc,
		cmdExecutor,
		dynamoDBSvc,
		mgr,
	)

	return []Controller{
		NewPingController(),
		NewSlackController(svc),
		NewEveController(svc),
		NewAuthController(mgr, dynamoDBSvc),
	}
}
