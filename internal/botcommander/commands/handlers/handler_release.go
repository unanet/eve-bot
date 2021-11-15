package handlers

import (
	"context"
	"fmt"
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/eve-bot/internal/service"
	"github.com/unanet/eve/pkg/eve"
	"strings"
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

// handleArtifactRelease handles the ReleaseCmd for releasing a namespace
func (h ReleaseHandler) handleArtifactRelease(ctx context.Context, cmd commands.EvebotCommand, timestamp string, dynamicOpts commands.CommandOptions) {

	releases, err := h.svc.EveAPI.Release(ctx, eve.Release{
		Type:     eve.ReleaseTypeArtifact,
		Artifact: dynamicOpts[params.ArtifactName].(string),
		Version:  dynamicOpts[params.ArtifactVersionName].(string),
		FromFeed: dynamicOpts[params.FromFeedName].(string),
		ToFeed:   dynamicOpts[params.ToFeedName].(string),
	})

	h.handleReleaseOutput(ctx, cmd, timestamp, releases, err)
}

// handleNamespaceRelease handles the ReleaseCmd for releasing a namespace
func (h ReleaseHandler) handleNamespaceRelease(ctx context.Context, cmd commands.EvebotCommand, timestamp string, dynamicOpts commands.CommandOptions) {

	releases, err := h.svc.EveAPI.Release(ctx, eve.Release{
		Type:        eve.ReleaseTypeNamespace,
		Namespace:   dynamicOpts[params.NamespaceName].(string),
		Environment: dynamicOpts[params.EnvironmentName].(string),
		FromFeed:    dynamicOpts[params.FromFeedName].(string),
		ToFeed:      dynamicOpts[params.ToFeedName].(string),
	})

	h.handleReleaseOutput(ctx, cmd, timestamp, releases, err)
}

// handleReleaseOutput handles the output from EveAPI into a consistent readable message
func (h ReleaseHandler) handleReleaseOutput(ctx context.Context, cmd commands.EvebotCommand, timestamp string, releases []eve.Release, err error) {
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("failed release: %s", err.Error()), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	var messages []string
	for _, rel := range releases {
		messages = append(messages, eveapi.ChatMessage(rel))
	}

	h.svc.ChatService.ReleaseResultsMessageThread(ctx, strings.Join(messages, "\n\n\n"), cmd.Info().User, cmd.Info().Channel, timestamp)
}