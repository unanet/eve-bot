package handlers

import (
	"context"
	"fmt"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
	"strings"

	"github.com/unanet/eve-bot/internal/service"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/eve/pkg/eve"
)

// ReleaseHandler is the handler for the ReleaseCmd
type ReleaseHandler struct {
	svc *service.Provider
}

// NewReleaseHandler creates a ReleaseHandler
func NewReleaseHandler(svc *service.Provider) CommandHandler {
	return ReleaseHandler{svc: svc}
}

// Handle handles the ReleaseCmd
func (h ReleaseHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

	dynamicOpts := cmd.Options()

	switch dynamicOpts["commandType"] {
		case "artifact":
			h.handleArtifactRelease(ctx, cmd, timestamp, dynamicOpts)
		case "namespace":
			h.handleNamespaceRelease(ctx, cmd, timestamp, dynamicOpts)
		default:
			h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("unable to handle release type of (%s)", dynamicOpts["CommandType"]), cmd.Info().User, cmd.Info().Channel, timestamp)
	}
}

// HandleNamespaceRelease handles the ReleaseCmd for releasing a namespace
func (h ReleaseHandler) handleArtifactRelease(ctx context.Context, cmd commands.EvebotCommand, timestamp string, dynamicOpts commands.CommandOptions) {

	release, err := h.svc.EveAPI.Release(ctx, eve.Release{
		Artifact: dynamicOpts[params.ArtifactName].(string),
		Version:  dynamicOpts[params.ArtifactVersionName].(string),
		FromFeed: dynamicOpts[params.FromFeedName].(string),
		ToFeed:   dynamicOpts[params.ToFeedName].(string),
	})
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("failed release: %s", err.Error()), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	h.svc.ChatService.ReleaseResultsMessageThread(ctx, eveapi.ChatMessage(release), cmd.Info().User, cmd.Info().Channel, timestamp)
}

func (h ReleaseHandler) artifactsToReleaseForNamespace(ctx context.Context, cmd commands.EvebotCommand, ts string) ([]eve.Service, error) {
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

// HandleNamespaceRelease handles the ReleaseCmd for releasing a namespace
func (h ReleaseHandler) handleNamespaceRelease(ctx context.Context, cmd commands.EvebotCommand, timestamp string, dynamicOpts commands.CommandOptions) {

	services, err := h.artifactsToReleaseForNamespace(ctx, cmd, timestamp)
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
		release, err := h.svc.EveAPI.Release(ctx, eve.Release{
			Artifact: svc.ArtifactName,
			Version:  svc.DeployedVersion,
			FromFeed: dynamicOpts[params.FromFeedName].(string),
			ToFeed:   dynamicOpts[params.ToFeedName].(string),
		})

		if err != nil {
			failedArtifacts = append(failedArtifacts, fmt.Sprintf("%s:%s", svc.ArtifactName, svc.DeployedVersion))
			errs = append(errs, err)
			continue
		}

		builder.WriteString(eveapi.ChatMessage(release) + "\n\n\n") // triple new line is intentional due to formatting

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

	h.svc.ChatService.ShowResultsMessageThread(ctx, builder.String(), cmd.Info().User, cmd.Info().Channel, timestamp)
}