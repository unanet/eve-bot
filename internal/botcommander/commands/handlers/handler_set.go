package handlers

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"strings"

	"gitlab.unanet.io/devops/eve/pkg/eve"
	"gitlab.unanet.io/devops/eve/pkg/log"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
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
	ns, err := resolveNamespace(ctx, h.eveAPIClient, cmd)
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, err.Error(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
		return
	}

	svcs, err := h.eveAPIClient.GetServicesByNamespace(ctx, ns.Name)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp, err)
		return
	}
	if svcs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no services", cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
		return
	}

	// Service was supplied (we are setting the resource at the service level)
	if _, ok := cmd.APIOptions()[params.ServiceName].(string); ok {
		var svc eveapimodels.EveService
		for _, s := range svcs {
			if strings.ToLower(s.Name) == strings.ToLower(cmd.APIOptions()[params.ServiceName].(string)) {
				svc = mapToEveService(s)
				break
			}
		}
		if svc.ID == 0 {
			h.chatSvc.UserNotificationThread(ctx, "invalid service", cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
			return
		}

		switch cmd.APIOptions()["resource"] {
		case resources.MetadataName:
			h.setSvcMetadata(ctx, cmd, &timestamp, svc)
			return
		case resources.VersionName:
			h.setSvcVersion(ctx, cmd, &timestamp, svc)
			return
		}
	}

	// setting the resource at the namespace level
	switch cmd.APIOptions()["resource"] {
	case resources.MetadataName:
		h.chatSvc.UserNotificationThread(ctx, "cannot set namespace metadata", cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
	case resources.VersionName:
		h.setNamespaceVersion(ctx, cmd, &timestamp, ns)
	}
}

func (h SetHandler) setSvcMetadata(ctx context.Context, cmd commands.EvebotCommand, ts *string, svc eveapimodels.EveService) {
	var metadataMap params.MetadataMap
	var metaDataOK bool
	log.Logger.Debug("cmd.Api.Options", zap.Any("api_opts", cmd.APIOptions()))
	if metadataMap, metaDataOK = cmd.APIOptions()[params.MetadataName].(params.MetadataMap); !metaDataOK {

		if metadataMap == nil {
			log.Logger.Debug("cmd.Api.Options.metadatamap nil")
		} else {
			log.Logger.Debug("cmd.Api.Options.metadatamap nil", zap.Any("metadata_map", metadataMap))
		}

		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, fmt.Errorf("invalid metadata map"))
		return
	}

	md, err := h.eveAPIClient.SetServiceMetadata(ctx, metadataMap, svc.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, err)
		return
	}

	h.chatSvc.UserNotificationThread(ctx, md.ToString(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
}

func (h SetHandler) setSvcVersion(ctx context.Context, cmd commands.EvebotCommand, ts *string, svc eveapimodels.EveService) {
	var version string
	var versionOK bool
	if version, versionOK = cmd.APIOptions()[params.VersionName].(string); !versionOK {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, fmt.Errorf("invalid version"))
		return
	}
	updatedSvc, err := h.eveAPIClient.SetServiceVersion(ctx, version, svc.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, err)
		return
	}
	h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("%s version set to %s", updatedSvc.Name, updatedSvc.OverrideVersion), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
}

func (h SetHandler) setNamespaceVersion(ctx context.Context, cmd commands.EvebotCommand, ts *string, ns eve.Namespace) {
	var version string
	var versionOK bool
	if version, versionOK = cmd.APIOptions()[params.VersionName].(string); !versionOK {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, fmt.Errorf("invalid version"))
		return
	}
	updatedNS, err := h.eveAPIClient.SetNamespaceVersion(ctx, version, ns.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, err)
		return
	}
	h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("%s version set to %s", updatedNS.Name, updatedNS.RequestedVersion), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
}
