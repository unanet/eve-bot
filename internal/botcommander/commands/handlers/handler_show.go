package handlers

import (
	"context"
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

type ShowHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewShowHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return ShowHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h ShowHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	//chatUser, err := h.chatSvc.GetUser(ctx, cmd.User())
	//if err != nil {
	//	h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, err)
	//	return
	//}

	cmdAPIOpts := cmd.APIOptions()

	h.chatSvc.ShowResultsMessageThread(ctx, fmt.Sprintf("resource: %s", cmdAPIOpts["resource"]), cmd.User(), cmd.Channel(), timestamp)

	//deployOpts := eveapi.DeploymentPlanOptions{
	//	Artifacts:        commands.ExtractServiceArtifactsOpt(cmdAPIOpts),
	//	ForceDeploy:      commands.ExtractForceDeployOpt(cmdAPIOpts),
	//	User:             chatUser.Name,
	//	DryRun:           commands.ExtractDryrunOpt(cmdAPIOpts),
	//	Environment:      commands.ExtractEnvironmentOpt(cmdAPIOpts),
	//	NamespaceAliases: commands.ExtractNSOpt(cmdAPIOpts),
	//	Messages:         nil,
	//	Type:             "application",
	//}

	//resp, err := h.eveAPIClient.Deploy(ctx, deployOpts, cmd.User(), cmd.Channel(), timestamp)
	//if err != nil && len(err.Error()) > 0 {
	//	h.chatSvc.DeploymentNotificationThread(ctx, err.Error(), cmd.User(), cmd.Channel(), timestamp)
	//	return
	//}
	//if resp == nil {
	//	h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, errInvalidApiResp)
	//	return
	//}
	//if len(resp.Messages) > 0 {
	//	h.chatSvc.UserNotificationThread(ctx, strings.Join(resp.Messages, ","), cmd.User(), cmd.Channel(), timestamp)
	//	return
	//}
	return

}
