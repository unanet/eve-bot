package handlers

import (
	"context"

	"github.com/unanet/eve-bot/internal/service"

	"github.com/unanet/eve-bot/internal/botcommander/args"
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve/pkg/eve"
)

// DeployHandler is the handler for the DeployCmd
type DeployHandler struct {
	svc *service.Provider
}

// NewDeployHandler creates a DeployHandler
func NewDeployHandler(svc *service.Provider) CommandHandler {
	return DeployHandler{svc: svc}
}

// Handle handles the DeployCmd
func (h DeployHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	chatUser, err := h.svc.ChatService.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	cmdAPIOpts := cmd.Options()

	deployHandler(ctx, h.svc.EveAPI, h.svc.ChatService, cmd, timestamp, eve.DeploymentPlanOptions{
		Artifacts:        commands.ExtractArtifactsDefinition(args.ServicesName, cmdAPIOpts),
		ForceDeploy:      commands.ExtractBoolOpt(args.ForceDeployName, cmdAPIOpts),
		User:             chatUser.Name,
		DryRun:           commands.ExtractBoolOpt(args.DryrunName, cmdAPIOpts),
		Environment:      commands.ExtractStringOpt(params.EnvironmentName, cmdAPIOpts),
		NamespaceAliases: commands.ExtractStringListOpt(params.NamespaceName, cmdAPIOpts),
		Type:             eve.DeploymentPlanTypeApplication,
	})
}
