package handlers

import (
	"context"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"

	"github.com/unanet/eve-bot/internal/botcommander/args"
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve/pkg/eve"
)

// DeployHandler is the handler for the DeployCmd
type DeployHandler struct {
	eveAPIClient interfaces.EveAPI
	chatSvc      interfaces.ChatProvider
}

// NewDeployHandler creates a DeployHandler
func NewDeployHandler(eveAPIClient interfaces.EveAPI, chatSvc interfaces.ChatProvider) CommandHandler {
	return DeployHandler{
		eveAPIClient: eveAPIClient,
		chatSvc:      chatSvc,
	}
}

// Handle handles the DeployCmd
func (h DeployHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	chatUser, err := h.chatSvc.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	cmdAPIOpts := cmd.Options()

	deployHandler(ctx, h.eveAPIClient, h.chatSvc, cmd, timestamp, eve.DeploymentPlanOptions{
		Artifacts:        commands.ExtractArtifactsDefinition(args.ServicesName, cmdAPIOpts),
		ForceDeploy:      commands.ExtractBoolOpt(args.ForceDeployName, cmdAPIOpts),
		User:             chatUser.Name,
		DryRun:           commands.ExtractBoolOpt(args.DryrunName, cmdAPIOpts),
		Environment:      commands.ExtractStringOpt(params.EnvironmentName, cmdAPIOpts),
		NamespaceAliases: commands.ExtractStringListOpt(params.NamespaceName, cmdAPIOpts),
		Type:             eve.DeploymentPlanTypeApplication,
	})
}
