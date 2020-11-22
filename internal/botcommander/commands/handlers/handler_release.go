package handlers

import (
	"context"
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

// ReleaseHandler is the handler for the ReleaseCmd
type ReleaseHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

// NewReleaseHandler creates a ReleaseHandler
func NewReleaseHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return ReleaseHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

// Handle handles the ReleaseCmd
func (h ReleaseHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

	dynamicOpts := cmd.Options()

	release, err := h.eveAPIClient.Release(ctx, eve.Release{
		Artifact: dynamicOpts[params.ArtifactName].(string),
		Version:  dynamicOpts[params.ArtifactVersionName].(string),
		FromFeed: dynamicOpts[params.FromFeedName].(string),
		ToFeed:   dynamicOpts[params.ToFeedName].(string),
	})
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("failed release: %s", err.Error()), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	h.chatSvc.ReleaseResultsMessageThread(ctx, eveapi.ToChatMessage(release), cmd.Info().User, cmd.Info().Channel, timestamp)
}
