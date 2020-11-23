package handlers

import (
	"context"
	"fmt"

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

// Handle handles the RestartCmd
func (h RestartHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	ns, svc := resolveServiceNamespace(ctx, h.eveAPIClient, h.chatSvc, cmd, &timestamp)
	if ns == nil || svc == nil {
		h.chatSvc.UserNotificationThread(ctx, "failed to resolve the restart command service and namespace params", cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	chatUser, err := h.chatSvc.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	if len(svc.DeployedVersion) == 0 {
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("can't restart deployed_version: %s (empty)", svc.DeployedVersion), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	deployOpts := eve.DeploymentPlanOptions{
		Artifacts: eve.ArtifactDefinitions{
			&eve.ArtifactDefinition{
				Name:             commands.ExtractStringOpt(params.ServiceName, cmd.Options()),
				RequestedVersion: svc.DeployedVersion,
			},
		},
		ForceDeploy:      true,
		User:             chatUser.Name,
		DryRun:           false,
		Environment:      commands.ExtractStringOpt(params.EnvironmentName, cmd.Options()),
		NamespaceAliases: commands.ExtractStringListOpt(params.NamespaceName, cmd.Options()),
		Type:             eve.DeploymentPlanTypeRestart,
	}

	deployHandler(ctx, h.eveAPIClient, h.chatSvc, cmd, timestamp, deployOpts)

}
