package handlers

import (
	"context"
	"strings"

	"github.com/unanet/eve-bot/internal/service"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve/pkg/eve"
)

// RunHandler is the handler for the RunCmd
type RunHandler struct {
	svc *service.Provider
}

// NewRunHandler creates a RunHandler
func NewRunHandler(svc *service.Provider) CommandHandler {
	return RunHandler{svc: svc}
}

// Handle handles the RunCmd
func (h RunHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	chatUser, err := h.svc.ChatService.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	cmdAPIOpts := cmd.Options()

	// TODO: Get this out of the handler and into the command
	//  ideally we resolve this data in the command, before passing to the handler
	var aDefs eve.ArtifactDefinitions

	if job, ok := cmdAPIOpts[params.JobName].(string); ok {
		aDef := &eve.ArtifactDefinition{}
		if strings.Contains(job, ":") {
			kv := strings.Split(job, ":")
			aDef.Name = kv[0]
			aDef.RequestedVersion = kv[1]
		} else {
			aDef.Name = job
		}
		aDefs = append(aDefs, aDef)
	}

	deployHandler(ctx, h.svc.EveAPI, h.svc.ChatService, cmd, timestamp, eve.DeploymentPlanOptions{
		Artifacts:        aDefs,
		ForceDeploy:      true,
		User:             chatUser.Name,
		DryRun:           false,
		Environment:      commands.ExtractStringOpt(params.EnvironmentName, cmdAPIOpts),
		NamespaceAliases: commands.ExtractStringListOpt(params.NamespaceName, cmdAPIOpts),
		Messages:         nil,
		Type:             eve.DeploymentPlanTypeJob,
		Metadata:         commands.ExtractMetadataField(cmdAPIOpts),
	})
}
