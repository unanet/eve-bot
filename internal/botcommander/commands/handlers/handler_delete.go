package handlers

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
	"gitlab.unanet.io/devops/eve/pkg/log"
)

type DeleteHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewDeleteHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return DeleteHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h DeleteHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	nv, err := resolveNamespace(ctx, h.eveAPIClient, cmd)
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, err.Error(), cmd.User(), cmd.Channel(), timestamp)
		return
	}
	svcs, err := h.eveAPIClient.GetServicesByNamespace(ctx, nv.Name)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, err)
		return
	}
	if svcs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no services", cmd.User(), cmd.Channel(), timestamp)
		return
	}
	log.Logger.Debug("services", zap.Any("svcs", svcs))
	var svc eveapimodels.EveService
	for _, s := range svcs {
		if strings.ToLower(s.Name) == strings.ToLower(cmd.APIOptions()[params.ServiceName].(string)) {
			svc = mapToEveService(s)
			break
		}
	}
	log.Logger.Debug("service", zap.Any("svc", svc))
	var md params.MetadataMap

	for _, m := range cmd.APIOptions()[params.MetadataName].([]string) {
		log.Logger.Debug("DeleteServiceMetadata", zap.Any("svc", svc), zap.String("metadata_key", m))
		md, err = h.eveAPIClient.DeleteServiceMetadata(ctx, m, svc.ID)
		if err != nil {
			h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, err)
			return
		}
	}

	if md == nil {
		h.chatSvc.UserNotificationThread(ctx, "no metadata", cmd.User(), cmd.Channel(), timestamp)
		return
	}
	h.chatSvc.UserNotificationThread(ctx, md.ToString(), cmd.User(), cmd.Channel(), timestamp)
}
