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
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

type ShowHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewShowHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return ShowHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h ShowHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	switch cmd.DynamicOptions()["resource"] {
	case resources.EnvironmentName:
		h.showEnvironments(ctx, cmd, &timestamp)
	case resources.NamespaceName:
		h.showNamespaces(ctx, cmd, &timestamp)
	case resources.ServiceName:
		h.showServices(ctx, cmd, &timestamp)
	case resources.MetadataName:
		h.showMetadata(ctx, cmd, &timestamp)
	default:
		h.chatSvc.UserNotificationThread(ctx, "invalid show command", cmd.ChatInfo().User, cmd.ChatInfo().Channel, timestamp)
	}
}

func (h ShowHandler) showEnvironments(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	envs, err := h.eveAPIClient.GetEnvironments(ctx)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, err)
		return
	}
	if envs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no environments", cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, envs.ToChatMessage(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
}

func (h ShowHandler) showNamespaces(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, err := h.eveAPIClient.GetNamespacesByEnvironment(ctx, cmd.DynamicOptions()[params.EnvironmentName].(string))
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, err)
		return
	}
	if ns == nil {
		h.chatSvc.UserNotificationThread(ctx, "no namespaces", cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, ns.ToChatMessage(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
}

func (h ShowHandler) showServices(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	nv, err := resolveNamespace(ctx, h.eveAPIClient, cmd)
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, err.Error(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
		return
	}
	svcs, err := h.eveAPIClient.GetServicesByNamespace(ctx, nv.Name)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, err)
		return
	}
	if svcs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no services", cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, svcs.ToChatMessage(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
}

func (h ShowHandler) showMetadata(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	nv, err := resolveNamespace(ctx, h.eveAPIClient, cmd)
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, err.Error(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
		return
	}
	svcs, err := h.eveAPIClient.GetServicesByNamespace(ctx, nv.Name)
	log.Logger.Debug("services", zap.Any("svcs", svcs))
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, err)
		return
	}
	if svcs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no services", cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
		return
	}
	var requestedSvcName string
	var valid bool
	if requestedSvcName, valid = cmd.DynamicOptions()[params.ServiceName].(string); !valid {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, fmt.Errorf("invalid ServiceName Param"))
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
		h.chatSvc.UserNotificationThread(ctx, fmt.Sprintf("invalid requested service: %s", requestedSvcName), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
		return
	}
	log.Logger.Debug("server", zap.Any("svc", svc))
	fullSvc, err := h.eveAPIClient.GetServiceByID(ctx, svc.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts, err)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, fullSvc.MetadataToChatMessage(), cmd.ChatInfo().User, cmd.ChatInfo().Channel, *ts)
}
