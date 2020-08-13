package handlers

import (
	"context"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
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
	nv, err := resolveNamespace(ctx, h.eveAPIClient, cmd)
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}
	svcs, err := h.eveAPIClient.GetServicesByNamespace(ctx, nv.Name)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}
	if svcs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no services", cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}
	var requestedSvcName string
	var validSvc bool
	if requestedSvcName, validSvc = cmd.Options()[params.ServiceName].(string); !validSvc {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, fmt.Errorf("invalid ServiceName Param"))
		return
	}

	var svc eveapimodels.EveService
	for _, s := range svcs {
		if strings.ToLower(s.Name) == strings.ToLower(requestedSvcName) {
			svc = mapToEveService(s)
			break
		}
	}
	if svc.ID == 0 {
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("invalid requested service: %s", requestedSvcName), cmd.Info().User, cmd.Info().Channel, timestamp)
		return
	}

	switch cmd.Options()["resource"] {
	case resources.MetadataName:
		h.deleteMetadata(ctx, cmd, &timestamp, svc)
	case resources.VersionName:
		h.deleteVersion(ctx, cmd, &timestamp, svc)
	}
}

func (h DeleteHandler) deleteMetadata(ctx context.Context, cmd commands.EvebotCommand, ts *string, svc eveapimodels.EveService) {
	var requestedMetadata []string
	var validMetadata bool
	if requestedMetadata, validMetadata = cmd.Options()[params.MetadataName].([]string); validMetadata == false {
		log.Logger.Warn("troy debug invalid metadata",
			zap.Any("opts", cmd.Options()),
			zap.Strings("requestedMetadata", requestedMetadata),
			zap.Bool("validMetadata", validMetadata))
		h.chatSvc.UserNotificationThread(ctx, "invalid metadata", cmd.Info().User, cmd.Info().Channel, *ts)
		//h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, fmt.Errorf("invalid MetadataName Param"))
		return
	}
	if len(requestedMetadata) == 0 {
		h.chatSvc.UserNotificationThread(ctx, "you must supply 1 or more metadata keys", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	var md params.MetadataMap
	var err error
	for _, m := range requestedMetadata {
		md, err = h.eveAPIClient.DeleteServiceMetadata(ctx, m, svc.ID)
		if err != nil {
			h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
			return
		}
	}
	if md == nil {
		h.chatSvc.UserNotificationThread(ctx, "no metadata", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	h.chatSvc.UserNotificationThread(ctx, md.ToString(), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h DeleteHandler) deleteVersion(ctx context.Context, cmd commands.EvebotCommand, ts *string, svc eveapimodels.EveService) {
	updatedSvc, err := h.eveAPIClient.SetServiceVersion(ctx, "", svc.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("%s version deleted", updatedSvc.Name), cmd.Info().User, cmd.Info().Channel, *ts)
}
