package handlers

import (
	"context"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/resources"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
)

// ShowHandler is the handler for the ShowCmd
type ShowHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

// NewShowHandler creates a ShowHandler
func NewShowHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return ShowHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

// Handle handles the ShowCmd
func (h ShowHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	switch cmd.Options()["resource"] {
	case resources.EnvironmentName:
		h.showEnvironments(ctx, cmd, &timestamp)
	case resources.NamespaceName:
		h.showNamespaces(ctx, cmd, &timestamp)
	case resources.ServiceName:
		h.showServices(ctx, cmd, &timestamp)
	case resources.MetadataName:
		h.showMetadata(ctx, cmd, &timestamp)
	default:
		h.chatSvc.UserNotificationThread(ctx, "invalid show command", cmd.Info().User, cmd.Info().Channel, timestamp)
	}
}

func (h ShowHandler) showEnvironments(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	envs, err := h.eveAPIClient.GetEnvironments(ctx)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	if envs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no environments", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, envs.ToChatMessage(), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h ShowHandler) showNamespaces(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, err := h.eveAPIClient.GetNamespacesByEnvironment(ctx, cmd.Options()[params.EnvironmentName].(string))
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	if ns == nil {
		h.chatSvc.UserNotificationThread(ctx, "no namespaces", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, ns.ToChatMessage(), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h ShowHandler) showServices(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	nv, err := resolveNamespace(ctx, h.eveAPIClient, cmd)
	if err != nil {
		h.chatSvc.UserNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	svcs, err := h.eveAPIClient.GetServicesByNamespace(ctx, nv.Name)
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	if svcs == nil {
		h.chatSvc.UserNotificationThread(ctx, "no services", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	h.chatSvc.ShowResultsMessageThread(ctx, svcs.ToChatMessage(), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h ShowHandler) showMetadata(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, svc := resolveServiceNamespace(ctx, h.eveAPIClient, h.chatSvc, cmd, ts)
	if svc == nil || ns == nil {
		return
	}

	metadata, err := h.eveAPIClient.GetMetadata(ctx, metaDataServiceKey(svc.Name, ns.Name))
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}

	h.chatSvc.ShowResultsMessageThread(ctx, eveapimodels.MetaData{Input: metadata}.ToChatMessage(), cmd.Info().User, cmd.Info().Channel, *ts)
}
