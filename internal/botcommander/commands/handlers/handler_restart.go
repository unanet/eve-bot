package handlers

import (
	"context"
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
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

	ns, err := resolveNamespace(ctx, h.eveAPIClient, cmd)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	chatUser, err := h.chatSvc.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}

	var artifactDef eveapimodels.ArtifactDefinitions

	artifactsRequested := commands.ExtractArtifactsDefinition(args.ServicesName, cmd.Options())

	svcs, err := h.eveAPIClient.GetServicesByNamespace(ctx, ns.Name)
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("(GetServicesByNamespace) failed to find service: %s", err.Error()), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	if len(artifactsRequested) == 0 {
		// No specific services were requests, so we are "restarting" (redeploying the same version) of all services in the namesace
		artifactDef = servicesToArtifactDef(svcs)
	} else {
		// The user has requested specific services
		// lets map them first and then "find" them
		// ... a map here is faster than the nested for loops required
		svcMap := toServicesMap(svcs)
		var currServicesRequest []eve.Service
		for _, artifactReq := range artifactsRequested {
			currServicesRequest = append(currServicesRequest, svcMap[artifactReq.Name])
		}
		artifactDef = servicesToArtifactDef(currServicesRequest)
	}

	deployOpts := eveapimodels.DeploymentPlanOptions{
		Artifacts:        artifactDef,
		ForceDeploy:      true,
		User:             chatUser.Name,
		DryRun:           false,
		Environment:      commands.ExtractStringOpt(params.EnvironmentName, cmd.Options()),
		NamespaceAliases: commands.ExtractStringListOpt(params.NamespaceName, cmd.Options()),
		Type:             "application",
	}

	deployHandler(ctx, h.eveAPIClient, h.chatSvc, cmd, timestamp, deployOpts)

}

// we convert the slice of services to map["service_name"] = Service
func toServicesMap(svcs eveapimodels.Services) map[string]eve.Service {
	result := make(map[string]eve.Service)
	for _, svc := range svcs {
		result[svc.Name] = svc
	}
	return result
}

func servicesToArtifactDef(svcs eveapimodels.Services) eveapimodels.ArtifactDefinitions {
	var result eveapimodels.ArtifactDefinitions
	for _, svc := range svcs {
		def := &eveapimodels.ArtifactDefinition{
			Name:             svc.Name,
			RequestedVersion: svc.DeployedVersion,
		}
		result = append(result, def)
	}
	return result
}
