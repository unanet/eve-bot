package handlers

import (
	"context"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

func deployHandler(
	ctx context.Context,
	eveAPIClient eveapi.Client,
	chatSvc chatservice.Provider,
	cmd commands.EvebotCommand,
	timestamp string,
	deployOpts eveapimodels.DeploymentPlanOptions) {

	resp, err := eveAPIClient.Deploy(ctx, deployOpts, cmd.Info().User, cmd.Info().Channel, timestamp)
	if err != nil && len(err.Error()) > 0 {
		chatSvc.DeploymentNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}
	if resp == nil {
		chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, errInvalidAPIResp)
		return
	}
	if len(resp.Messages) > 0 {
		chatSvc.UserNotificationThread(ctx, strings.Join(resp.Messages, ","), cmd.Info().User, cmd.Info().Channel, timestamp)
	}
}
