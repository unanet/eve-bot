package handlers

import (
	"context"
	"fmt"

	"github.com/unanet/eve-bot/internal/service"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve/pkg/eve"
)

// ReleaseHandler is the handler for the ReleaseCmd
type RestartHandler struct {
	svc *service.Provider
}

// NewReleaseHandler creates a ReleaseHandler
func NewRestartHandler(svc *service.Provider) CommandHandler {
	return RestartHandler{svc: svc}
}

// Handle handles the RestartCmd
func (h RestartHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	ns, svc := resolveServiceNamespace(ctx, h.svc.EveAPI, h.svc.ChatService, cmd, &timestamp)
	if ns == nil || svc == nil {
		h.svc.ChatService.UserNotificationThread(ctx, "failed to resolve the restart command service and namespace params", cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	chatUser, err := h.svc.ChatService.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	if len(svc.DeployedVersion) == 0 {
		h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("can't restart deployed_version: %s (empty)", svc.DeployedVersion), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	deployHandler(ctx, h.svc.EveAPI, h.svc.ChatService, cmd, timestamp, eve.DeploymentPlanOptions{
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
	})

}
