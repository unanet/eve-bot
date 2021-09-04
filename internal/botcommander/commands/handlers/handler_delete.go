package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/unanet/eve-bot/internal/service"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/botcommander/resources"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/eve/pkg/eve"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

// DeleteHandler is the handler for the DeleteCmd
type DeleteHandler struct {
	svc *service.Provider
}

// NewDeleteHandler creates a DeleteHandler
func NewDeleteHandler(svc *service.Provider) CommandHandler {
	return DeleteHandler{svc: svc}
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
		h.svc.ChatService.UserNotificationThread(ctx, "you must supply 1 or more metadata keys", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}

	ns, svc := resolveServiceNamespace(ctx, h.svc.EveAPI, h.svc.ChatService, cmd, ts)
	if svc == nil || ns == nil {
		return
	}

	mdKey := metaDataServiceKey(svc.Name, ns.Name)
	mdItem, err := h.svc.EveAPI.GetMetadata(ctx, mdKey)
	if err != nil {
		if resourceNotFoundError(err) {
			h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("no metadata found for: %s", mdKey), cmd.Info().User, cmd.Info().Channel, *ts)
			return
		}
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}

	var md eve.Metadata
	for _, m := range opts[params.MetadataName].([]string) {
		if isValidMetadata(m) {
			md, err = h.svc.EveAPI.DeleteMetadataKey(ctx, mdItem.ID, m)
			if err != nil {
				if resourceNotFoundError(err) {
					h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("failed to delete metadata key: %s", m), cmd.Info().User, cmd.Info().Channel, *ts)
					return
				}
				h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
				return
			}
		} else {
			h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("invalid metadata key: %s", m), cmd.Info().User, cmd.Info().Channel, *ts)
		}
	}

	h.svc.ChatService.ShowResultsMessageThread(ctx, eveapi.ChatMessage(md), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h DeleteHandler) deleteVersion(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, svc := resolveServiceNamespace(ctx, h.svc.EveAPI, h.svc.ChatService, cmd, ts)
	if svc == nil || ns == nil {
		return
	}

	updatedSvc, err := h.svc.EveAPI.SetServiceVersion(ctx, "", svc.ID)
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("%s version deleted", updatedSvc.Name), cmd.Info().User, cmd.Info().Channel, *ts)
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
