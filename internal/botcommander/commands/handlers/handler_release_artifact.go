package handlers

import (
	"context"
	"fmt"

	"github.com/unanet/eve-bot/internal/service"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/eve/pkg/eve"
)

// ReleaseArtifactHandler is the handler for the ReleaseArtifactCmd
type ReleaseArtifactHandler struct {
	svc *service.Provider
}

// NewReleaseArtifactHandler creates a ReleaseArtifactHandler
func NewReleaseArtifactHandler(svc *service.Provider) CommandHandler {
	return ReleaseArtifactHandler{svc: svc}
}

// Handle handles the ReleaseArtifactCmd
func (h ReleaseArtifactHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

	dynamicOpts := cmd.Options()

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
