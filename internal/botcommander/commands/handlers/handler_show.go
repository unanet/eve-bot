package handlers

import (
	"context"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
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
	switch cmd.APIOptions()["resource"] {
	case resources.EnvironmentName:
		h.showEnvironments(ctx, cmd, &timestamp)
	case resources.NamespaceName:
		h.showNamespaces(ctx, cmd, &timestamp)
	case resources.ServiceName:
		h.showServices(ctx, cmd, &timestamp)
	default:
		h.chatSvc.UserNotificationThread(ctx, "invalid show command", cmd.User(), cmd.Channel(), timestamp)
	}
}

func (h ShowHandler) showEnvironments(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	envs, err := h.eveAPIClient.GetEnvironments(ctx)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), *ts, err)
		return
	}
	if envs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no environments", cmd.User(), cmd.Channel(), *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, envs.ToChatMessage(), cmd.User(), cmd.Channel(), *ts)
}

func (h ShowHandler) showNamespaces(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, err := h.eveAPIClient.GetNamespacesByEnvironment(ctx, cmd.APIOptions()[params.EnvironmentName].(string))
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), *ts, err)
		return
	}
	if ns == nil {
		h.chatSvc.UserNotificationThread(ctx, "no namespaces", cmd.User(), cmd.Channel(), *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, ns.ToChatMessage(), cmd.User(), cmd.Channel(), *ts)
}

func (h ShowHandler) showServices(ctx context.Context, cmd commands.EvebotCommand, ts *string) {

	var nv eveapimodels.Namespace

	// Gotta get the namespaces first, since we are working with the Alias, and not the Name/ID
	namespaces, err := h.eveAPIClient.GetNamespacesByEnvironment(ctx, cmd.APIOptions()[params.EnvironmentName].(string))

	for _, v := range namespaces {
		if strings.ToLower(v.Alias) == strings.ToLower(cmd.APIOptions()[params.NamespaceName].(string)) {
			nv = v
			break
		}
	}

	if nv.ID == 0 {
		h.chatSvc.UserNotificationThread(ctx, "invalid namespace request", cmd.User(), cmd.Channel(), *ts)
		return
	}

	svc, err := h.eveAPIClient.GetServicesByNamespace(ctx, nv.Name)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), *ts, err)
		return
	}
	if svc == nil {
		h.chatSvc.UserNotificationThread(ctx, "no namespaces", cmd.User(), cmd.Channel(), *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, svc.ToChatMessage(), cmd.User(), cmd.Channel(), *ts)
}
