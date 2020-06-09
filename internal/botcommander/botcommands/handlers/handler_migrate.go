package handlers

import (
	"context"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
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

func (h MigrateHandler) Handle(ctx context.Context, cmd botcommands.EvebotCommand, timestamp string) {
	chatUser, err := h.chatSvc.GetUser(ctx, cmd.User())
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, err)
		return
	}

	cmdAPIOpts := cmd.APIOptions()

	deployOpts := eveapi.DeploymentPlanOptions{
		Artifacts:        botcommands.ExtractArtifactsOpt(cmdAPIOpts),
		ForceDeploy:      botcommands.ExtractForceDeployOpt(cmdAPIOpts),
		User:             chatUser.Name,
		DryRun:           botcommands.ExtractDryrunOpt(cmdAPIOpts),
		Environment:      botcommands.ExtractEnvironmentOpt(cmdAPIOpts),
		NamespaceAliases: botcommands.ExtractNSOpt(cmdAPIOpts),
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
		return
	}
	return

}
