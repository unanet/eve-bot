package handlers

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type ReleaseHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewReleaseHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return ReleaseHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h ReleaseHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {

	dynamicOpts := cmd.DynamicOptions()

	resp, err := h.eveAPIClient.Release(ctx, eve.Release{
		Artifact: dynamicOpts[params.ArtifactName].(string),
		Version:  dynamicOpts[params.ArtifactVersionName].(string),
		FromFeed: dynamicOpts[params.FromFeedName].(string),
		ToFeed:   dynamicOpts[params.ToFeedName].(string),
	})
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("failed release: %s", err.Error()), cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
		return
	}

	log.Logger.Debug("release response", zap.String("message", resp.Message))
	h.chatSvc.ReleaseResultsMessageThread(ctx, toChatMessage(resp), cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
}

func toChatMessage(resp eve.Release) string {
	return fmt.Sprintf("Artifact: `%s`\nVersion: `%s`\nFrom: `%s`\nTo: `%s`", resp.Artifact, resp.Version, resp.FromFeed, resp.ToFeed)
}
