package handlers

import (
	"context"
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
	"gitlab.unanet.io/devops/eve/pkg/eve"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
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
	nv, err := h.resolveNamespace(ctx, cmd)
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

	//fullSvc, err := h.eveAPIClient.GetServiceByID(ctx, svc.ID)
	//if err != nil {
	//	h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), *ts, err)
	//	return
	//}

	metadata := cmd.APIOptions()[params.MetadataName].(params.MetadataMap)

	//h.chatSvc.UserNotificationThread(ctx, metadata.ToString(), cmd.User(), cmd.Channel(), timestamp)

	md, err := h.eveAPIClient.SetServiceMetadata(ctx, metadata, svc.ID)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, err)
		return
	}

	h.chatSvc.UserNotificationThread(ctx, md.ToString(), cmd.User(), cmd.Channel(), timestamp)

}

func (h SetHandler) resolveNamespace(ctx context.Context, cmd commands.EvebotCommand) (eve.Namespace, error) {
	var nv eve.Namespace

	// Gotta get the namespaces first, since we are working with the Alias, and not the Name/ID
	namespaces, err := h.eveAPIClient.GetNamespacesByEnvironment(ctx, cmd.APIOptions()[params.EnvironmentName].(string))

	if err != nil {
		return nv, err
	}

	for _, v := range namespaces {
		if strings.ToLower(v.Alias) == strings.ToLower(cmd.APIOptions()[params.NamespaceName].(string)) {
			nv = v
			break
		}
	}

	if nv.ID == 0 {
		return nv, fmt.Errorf("invalid namespace")
	}
	return nv, nil
}
