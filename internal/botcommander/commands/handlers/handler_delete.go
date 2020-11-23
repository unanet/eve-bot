package handlers

import (
	"context"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// DeleteHandler is the handler for the DeleteCmd
type DeleteHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

// NewDeleteHandler creates a DeleteHandler
func NewDeleteHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return DeleteHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

// Handle handles the DeleteCmd
func (h DeleteHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	switch cmd.Options()["resource"] {
	case resources.MetadataName:
		h.deleteMetadata(ctx, cmd, &timestamp)
	case resources.VersionName:
		h.deleteVersion(ctx, cmd, &timestamp)
	}
}

func (h DeleteHandler) deleteMetadata(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	opts := cmd.Options()
	if len(opts[params.MetadataName].([]string)) == 0 {
		h.chatSvc.UserNotificationThread(ctx, "you must supply 1 or more metadata keys", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}

	ns, svc := resolveServiceNamespace(ctx, h.eveAPIClient, h.chatSvc, cmd, ts)
	if svc == nil || ns == nil {
		return
	}

	mdKey := metaDataServiceKey(svc.Name, ns.Name)
	mdItem, err := h.eveAPIClient.GetMetadata(ctx, mdKey)
	if err != nil {
		if resourceNotFoundError(err) {
			h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("no metadata found for: %s", mdKey), cmd.Info().User, cmd.Info().Channel, *ts)
			return
		}
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}

	var md eve.Metadata
	for _, m := range opts[params.MetadataName].([]string) {
		if isValidMetadata(m) {
			md, err = h.eveAPIClient.DeleteMetadataKey(ctx, mdItem.ID, m)
			if err != nil {
				if resourceNotFoundError(err) {
					h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("failed to delete metadata key: %s", m), cmd.Info().User, cmd.Info().Channel, *ts)
					return
				}
				h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
				return
			}
		} else {
			h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("invalid metadata key: %s", m), cmd.Info().User, cmd.Info().Channel, *ts)
		}
	}

	h.chatSvc.ShowResultsMessageThread(ctx, eveapi.ToChatMessage(md), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h DeleteHandler) deleteVersion(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, svc := resolveServiceNamespace(ctx, h.eveAPIClient, h.chatSvc, cmd, ts)
	if svc == nil || ns == nil {
		return
	}

	updatedSvc, err := h.eveAPIClient.SetServiceVersion(ctx, "", svc.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("%s version deleted", updatedSvc.Name), cmd.Info().User, cmd.Info().Channel, *ts)
}

func isValidMetadata(key string) bool {
	// Guard against the user sending key=value
	// we only want to send the key to the API
	metadatakey := key
	if strings.Contains(key, "=") {
		metadatakey = strings.Split(key, "=")[0]
	}
	if strings.Contains(metadatakey, "/") {
		log.Logger.Warn("metadata key contains slash", zap.String("metadatakey", metadatakey))
		return false
	}
	return true
}
