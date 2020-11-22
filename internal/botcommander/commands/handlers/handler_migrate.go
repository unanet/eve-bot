package handlers

import (
	"context"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/eve"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

// MigrateHandler is the handler for the MigrateCmd
type MigrateHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

// NewMigrateHandler creates a MigrateHandler
func NewMigrateHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return MigrateHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

// Handle handles the MigrateCmd
func (h MigrateHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	chatUser, err := h.chatSvc.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	cmdAPIOpts := cmd.Options()

	deployOpts := eve.DeploymentPlanOptions{
		Artifacts:        commands.ExtractArtifactsDefinition(args.DatabasesName, cmdAPIOpts),
		ForceDeploy:      commands.ExtractBoolOpt(args.ForceDeployName, cmdAPIOpts),
		User:             chatUser.Name,
		DryRun:           commands.ExtractBoolOpt(args.DryrunName, cmdAPIOpts),
		Environment:      commands.ExtractStringOpt(params.EnvironmentName, cmdAPIOpts),
		NamespaceAliases: commands.ExtractStringListOpt(params.NamespaceName, cmdAPIOpts),
		Messages:         nil,
		Type:             eve.DeploymentPlanTypeMigration,
	}

	resp, err := h.eveAPIClient.Deploy(ctx, deployOpts, cmd.Info().User, cmd.Info().Channel, timestamp)
	if err != nil && len(err.Error()) > 0 {
		h.chatSvc.DeploymentNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}
	if resp == nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, errInvalidAPIResp)
		return
	}
	if len(resp.Messages) > 0 {
		h.chatSvc.UserNotificationThread(ctx, strings.Join(resp.Messages, ","), cmd.Info().User, cmd.Info().Channel, timestamp)
	}
}
