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
)

// SetHandler is the handler for the SetCmd
type SetHandler struct {
	svc *service.Provider
}

// We use this stacking order as the default for all user (eve-bot) metadata
const stackingOrder = 400

// NewSetHandler creates a SetHandler
func NewSetHandler(svc *service.Provider) CommandHandler {
	return SetHandler{svc: svc}
}

// Handle handles the SetCmd
func (h SetHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	ns, err := resolveNamespace(ctx, h.svc.EveAPI, cmd)
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	svcs, err := h.svc.EveAPI.GetServicesByNamespace(ctx, ns.Name)
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}
	if svcs == nil {
		h.svc.ChatService.UserNotificationThread(ctx, "no services", cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	// Service was supplied (we are setting the resource at the service level)
	if _, ok := cmd.Options()[params.ServiceName].(string); ok {
		var svc eve.Service
		for _, s := range svcs {
			if strings.EqualFold(s.Name, cmd.Options()[params.ServiceName].(string)) {
				svc = s
				break
			}
		}
		if svc.ID == 0 {
			h.svc.ChatService.UserNotificationThread(ctx, "invalid service", cmd.Info().User, cmd.Info().Channel, timestamp)
			return
		}

		switch cmd.Options()["resource"] {
		case resources.MetadataName:
			h.setSvcMetadata(ctx, cmd, &timestamp, svc)
			return
		case resources.VersionName:
			h.setSvcVersion(ctx, cmd, &timestamp, svc)
			return
		}
	}

	// setting the resource at the namespace level
	switch cmd.Options()["resource"] {
	case resources.MetadataName:
		h.svc.ChatService.UserNotificationThread(ctx, "cannot set namespace metadata", cmd.Info().User, cmd.Info().Channel, timestamp)
	case resources.VersionName:
		h.setNamespaceVersion(ctx, cmd, &timestamp, ns)
	}
}

func (h SetHandler) setSvcMetadata(ctx context.Context, cmd commands.EvebotCommand, ts *string, svc eve.Service) {
	nv, err := resolveNamespace(ctx, h.svc.EveAPI, cmd)
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}

	md, err := h.svc.EveAPI.UpsertMergeMetadata(ctx, eve.Metadata{
		Description: metaDataServiceKey(svc.Name, nv.Name),
		Value:       commands.ExtractMetadataField(cmd.Options()),
	})
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, fmt.Errorf("failed to save metadata"))
		return
	}

	_, err = h.svc.EveAPI.UpsertMetadataServiceMap(ctx, eve.MetadataServiceMap{
		Description:   md.Description,
		MetadataID:    md.ID,
		ServiceID:     svc.ID,
		StackingOrder: stackingOrder,
	})
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, fmt.Errorf("failed to save metadata"))
		return
	}

	h.svc.ChatService.ShowResultsMessageThread(ctx, eveapi.ChatMessage(md), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h SetHandler) setSvcVersion(ctx context.Context, cmd commands.EvebotCommand, ts *string, svc eve.Service) {
	var version string
	var versionOK bool
	if version, versionOK = cmd.Options()[params.VersionName].(string); !versionOK {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, fmt.Errorf("invalid version"))
		return
	}
	updatedSvc, err := h.svc.EveAPI.SetServiceVersion(ctx, version, svc.ID)
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("%s version set to %s", updatedSvc.Name, updatedSvc.OverrideVersion), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h SetHandler) setNamespaceVersion(ctx context.Context, cmd commands.EvebotCommand, ts *string, ns eve.Namespace) {
	var version string
	var versionOK bool
	if version, versionOK = cmd.Options()[params.VersionName].(string); !versionOK {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, fmt.Errorf("invalid version"))
		return
	}
	updatedNS, err := h.svc.EveAPI.SetNamespaceVersion(ctx, version, ns.ID)
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("%s version set to %s", updatedNS.Name, updatedNS.RequestedVersion), cmd.Info().User, cmd.Info().Channel, *ts)
}
