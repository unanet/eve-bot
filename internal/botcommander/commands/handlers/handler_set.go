package handlers

import (
	"context"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
)

type SetHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewSetHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return SetHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h SetHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
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
	var svc eveapimodels.EveService
	for _, s := range svcs {
		if strings.ToLower(s.Name) == strings.ToLower(cmd.APIOptions()[params.ServiceName].(string)) {
			svc = mapToEveService(s)
			break
		}
	}
	if svc.ID == 0 {
		h.chatSvc.UserNotificationThread(ctx, "invalid service", cmd.User(), cmd.Channel(), timestamp)
		return
	}

	switch cmd.APIOptions()["resource"] {
	case resources.MetadataName:
		h.setMetadata(ctx, cmd, &timestamp, svc)
	case resources.VersionName:
		h.setVersion(ctx, cmd, &timestamp, svc)
	}

}

func (h SetHandler) setMetadata(ctx context.Context, cmd commands.EvebotCommand, ts *string, svc eveapimodels.EveService) {
	var metadataMap params.MetadataMap
	var metaDataOK bool
	if metadataMap, metaDataOK = cmd.APIOptions()[params.MetadataName].(params.MetadataMap); !metaDataOK {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), *ts, fmt.Errorf("invalid metadata map"))
		return
	}

	md, err := h.eveAPIClient.SetServiceMetadata(ctx, metadataMap, svc.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), *ts, err)
		return
	}

	h.chatSvc.UserNotificationThread(ctx, md.ToString(), cmd.User(), cmd.Channel(), *ts)
}

func (h SetHandler) setVersion(ctx context.Context, cmd commands.EvebotCommand, ts *string, svc eveapimodels.EveService) {
	var version string
	var versionOK bool
	if version, versionOK = cmd.APIOptions()[params.VersionName].(string); !versionOK {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), *ts, fmt.Errorf("invalid version"))
		return
	}
	updatedSvc, err := h.eveAPIClient.SetServiceVersion(ctx, version, svc.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), *ts, err)
		return
	}
	h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("%s version set to %s", updatedSvc.Name, updatedSvc.OverrideVersion), cmd.User(), cmd.Channel(), *ts)
}
