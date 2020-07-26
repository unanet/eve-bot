package handlers

import (
	"context"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

type MigrateHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewMigrateHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return MigrateHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h MigrateHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	chatUser, err := h.chatSvc.GetUser(ctx, cmd.User())
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, err)
		return
	}

	cmdAPIOpts := cmd.APIOptions()

	dbArtifacts := commands.ExtractDatabaseArtifactsOpt(cmdAPIOpts)

	log.Logger.Debug("migrate handler", zap.Any("opts", cmdAPIOpts), zap.Any("artifacts", dbArtifacts))

	deployOpts := eveapimodels.DeploymentPlanOptions{
		Artifacts:        dbArtifacts,
		ForceDeploy:      commands.ExtractForceDeployOpt(cmdAPIOpts),
		User:             chatUser.Name,
		DryRun:           commands.ExtractDryrunOpt(cmdAPIOpts),
		Environment:      commands.ExtractEnvironmentOpt(cmdAPIOpts),
		NamespaceAliases: commands.ExtractNSOpt(cmdAPIOpts),
		Messages:         nil,
		Type:             "migration",
	}

	resp, err := h.eveAPIClient.Deploy(ctx, deployOpts, cmd.User(), cmd.Channel(), timestamp)
	if err != nil && len(err.Error()) > 0 {
		h.chatSvc.DeploymentNotificationThread(ctx, err.Error(), cmd.User(), cmd.Channel(), timestamp)
		return
	}
	if resp == nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, errInvalidApiResp)
		return
	}
	if len(resp.Messages) > 0 {
		h.chatSvc.UserNotificationThread(ctx, strings.Join(resp.Messages, ","), cmd.User(), cmd.Channel(), timestamp)
	}
}
