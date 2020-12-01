package handlers

import (
	"context"
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

// RunHandler is the handler for the RunCmd
type RunHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

// NewRunHandler creates a RunHandler
func NewRunHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return RunHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

// Handle handles the RunCmd
func (h RunHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	chatUser, err := h.chatSvc.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	var metadataMap params.MetadataMap
	var metaDataOK bool

	if metadataMap, metaDataOK = cmd.Options()[params.MetadataName].(params.MetadataMap); !metaDataOK {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, fmt.Errorf("invalid metadata map"))
		return
	}

	var metadataField = make(eve.MetadataField)
	for k, v := range metadataMap {
		metadataField[k] = v
	}

	cmdAPIOpts := cmd.Options()

	deployHandler(ctx, h.eveAPIClient, h.chatSvc, cmd, timestamp, eve.DeploymentPlanOptions{
		Artifacts:        commands.ExtractArtifactsDefinition(params.JobName, cmdAPIOpts),
		ForceDeploy:      commands.ExtractBoolOpt(args.ForceDeployName, cmdAPIOpts),
		User:             chatUser.Name,
		DryRun:           commands.ExtractBoolOpt(args.DryrunName, cmdAPIOpts),
		Environment:      commands.ExtractStringOpt(params.EnvironmentName, cmdAPIOpts),
		NamespaceAliases: commands.ExtractStringListOpt(params.NamespaceName, cmdAPIOpts),
		Messages:         nil,
		Type:             eve.DeploymentPlanTypeJob,
		Metadata:         metadataField,
	})
}
