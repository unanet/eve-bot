package handlers

import (
	"context"
	"fmt"
	"strconv"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

// ReleaseHandler is the handler for the ReleaseCmd
type RestartHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

// NewReleaseHandler creates a ReleaseHandler
func NewRestartHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return RestartHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

//{
//"id": 83,
//"namespace_id": 28,
//"namespace_name": "una-qa-release",
//"artifact_id": 202,
//"artifact_name": "unanet-app",
//"override_version": "",
//"deployed_version": "20.7.0.5839",
//"created_at": "2020-06-22T18:16:58.349838Z",
//"updated_at": "2020-11-03T20:42:41.054868Z",
//"name": "unaneta",
//"sticky_sessions": true,
//"count": 1
//}

// Handle handles the ReleaseCmd
func (h RestartHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

	ns, err := resolveNamespace(ctx, h.eveAPIClient, cmd)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	svc, err := h.eveAPIClient.GetServiceByName(ctx, ns.Name, cmd.Options()[params.ServiceName].(string))
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("failed to find service: %s", err.Error()), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	chatUser, err := h.chatSvc.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	env, err := h.eveAPIClient.GetEnvironmentByID(ctx, strconv.Itoa(ns.EnvironmentID))
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	deployOpts := eveapimodels.DeploymentPlanOptions{
		Artifacts: eveapimodels.ArtifactDefinitions{
			&eveapimodels.ArtifactDefinition{
				Name:             svc.Name,
				RequestedVersion: svc.DeployedVersion,
			},
		},
		ForceDeploy:      true,
		User:             chatUser.Name,
		DryRun:           false,
		Environment:      env.Name,
		NamespaceAliases: []string{ns.Alias},
		Type:             "application",
	}

	deployHandler(ctx, h.eveAPIClient, h.chatSvc, cmd, timestamp, deployOpts)

}

func (h RestartHandler) toChatMessage(resp eve.Service) string {
	return fmt.Sprintf("Restarted Service: `%s`\nVersion: `%s`", resp.Name, resp.DeployedVersion)
}
