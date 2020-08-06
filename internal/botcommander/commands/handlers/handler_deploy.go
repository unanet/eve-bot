package handlers

import (
	"context"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
)

type DeployHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewDeployHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return DeployHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h DeployHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	chatUser, err := h.chatSvc.GetUser(ctx, cmd.ChatInfo().User)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp, err)
		return
	}

	cmdAPIOpts := cmd.APIOptions()

	deployOpts := eveapimodels.DeploymentPlanOptions{
		Artifacts:        commands.ExtractServiceArtifactsOpt(cmdAPIOpts),
		ForceDeploy:      commands.ExtractForceDeployOpt(cmdAPIOpts),
		User:             chatUser.Name,
		DryRun:           commands.ExtractDryrunOpt(cmdAPIOpts),
		Environment:      commands.ExtractEnvironmentOpt(cmdAPIOpts),
		NamespaceAliases: commands.ExtractNSOpt(cmdAPIOpts),
		Messages:         nil,
		Type:             "application",
	}

	resp, err := h.eveAPIClient.Deploy(ctx, deployOpts, cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
	if err != nil && len(err.Error()) > 0 {
		h.chatSvc.DeploymentNotificationThread(ctx, err.Error(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
		return
	}
	if resp == nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp, errInvalidApiResp)
		return
	}
	if len(resp.Messages) > 0 {
		h.chatSvc.UserNotificationThread(ctx, strings.Join(resp.Messages, ","), cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
	}
}
