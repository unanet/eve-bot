package handlers

import (
	"context"
	"fmt"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/service"
)

// ReleaseNamespaceHandler is the handler for the ReleaseNamespaceCmd
type ReleaseNamespaceHandler struct {
	svc *service.Provider
}

// NewReleaseNamespaceHandler creates a ReleaseNamespaceHandler
func NewReleaseNamespaceHandler(svc *service.Provider) CommandHandler {
	return ReleaseNamespaceHandler{svc: svc}
}

// Handle handles the ReleaseArtifactCmd
func (h ReleaseNamespaceHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

	dynamicOpts := cmd.Options()

	services, err := h.getArtifactsToReleaseForNamespace(ctx, cmd, timestamp)
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("failed release artifacts for namespace: %s", err.Error()), cmd.Info().User, cmd.Info().Channel, timestamp)
	}

	var errs []error

	var releasedArtifacts []string
	var failedArtifacts []string
	var builder strings.Builder

	namespaceName := dynamicOpts[params.NamespaceName].(string)
	fromFeed := dynamicOpts[params.FromFeedName].(string)
	toFeed := dynamicOpts[params.ToFeedName].(string)

	for _, svc := range services {
		if _, err := h.svc.EveAPI.Release(ctx, eve.Release{
			Artifact: svc.ArtifactName,
			Version:  svc.DeployedVersion,
			FromFeed: dynamicOpts[params.FromFeedName].(string),
			ToFeed:   dynamicOpts[params.ToFeedName].(string),
		}); err != nil {
			failedArtifacts = append(failedArtifacts, fmt.Sprintf("%s:%s", svc.ArtifactName, svc.DeployedVersion))
			errs = append(errs, err)
			continue
		}

		builder.WriteString(fmt.Sprintf("`%s` (%s)\n", svc.ArtifactName, svc.DeployedVersion))

		releasedArtifacts = append(releasedArtifacts, fmt.Sprintf("%s:%s", svc.ArtifactName, svc.DeployedVersion))
	}

	if len(errs) != 0 {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp,
			fmt.Errorf("failed to release the following artifact(s) for %s-%s\n\n%s",
				dynamicOpts[params.EnvironmentName].(string),
				namespaceName,
				strings.Join(failedArtifacts, "\n"),
			),
		)

		log.Logger.Error("error when releasing artifacts",
			zap.String("Environment", dynamicOpts[params.EnvironmentName].(string)),
			zap.String("Namespace", namespaceName),
			zap.Strings("FailedArtifacts", failedArtifacts),
			zap.Strings("ReleasedArtifacts", releasedArtifacts),
			zap.String("FromFeed", fromFeed),
			zap.String("ToFeed", toFeed),
			zap.Errors("Errs", errs),
		)

		return
	}

	toFeedMessage := ""
	if toFeed != "" {
		toFeedMessage = " to " + toFeed
	}

	h.svc.ChatService.ShowResultsMessageThread(ctx, fmt.Sprintf("The following artifacts have been released from %s%s:\n%s", fromFeed, toFeedMessage, builder.String()), cmd.Info().User, cmd.Info().Channel, timestamp)
}

func (h ReleaseNamespaceHandler) getArtifactsToReleaseForNamespace(ctx context.Context, cmd commands.EvebotCommand, ts string) ([]eve.Service, error) {
	nv, err := resolveNamespace(ctx, h.svc.EveAPI, cmd)
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, ts)
		return nil, fmt.Errorf("unable to get namespace")
	}

	svcs, err := h.svc.EveAPI.GetServicesByNamespace(ctx, nv.Name)
	if err != nil {
		if resourceNotFoundError(err) {
			h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("failed to get services in namespace: %s", nv.Name), cmd.Info().User, cmd.Info().Channel, ts)
			return nil, fmt.Errorf("unable to get services in namespace")
		}
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, ts, err)

		return nil, fmt.Errorf("error when getting services by namespace")
	}

	return svcs, nil
}