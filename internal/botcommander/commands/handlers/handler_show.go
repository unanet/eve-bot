package handlers

import (
	"context"
	"fmt"

	"github.com/unanet/eve-bot/internal/service"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/eve-bot/internal/botcommander/resources"
	"github.com/unanet/eve-bot/internal/eveapi"
	"github.com/unanet/go/pkg/errors"
)

// ShowHandler is the handler for the ShowCmd
type ShowHandler struct {
	svc *service.Provider
}

// NewShowHandler creates a ShowHandler
func NewShowHandler(svc *service.Provider) CommandHandler {
	return ShowHandler{svc: svc}
}

// Handle handles the ShowCmd
func (h ShowHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	switch cmd.Options()["resource"] {
	case resources.JobName, "jobs":
		h.showJobs(ctx, cmd, &timestamp)
	case resources.EnvironmentName:
		h.showEnvironments(ctx, cmd, &timestamp)
	case resources.NamespaceName:
		h.showNamespaces(ctx, cmd, &timestamp)
	case resources.ServiceName:
		h.showServices(ctx, cmd, &timestamp)
	case resources.MetadataName:
		h.showMetadata(ctx, cmd, &timestamp)
	default:
		h.svc.ChatService.UserNotificationThread(ctx, "invalid show command", cmd.Info().User, cmd.Info().Channel, timestamp)
	}
}

func (h ShowHandler) showEnvironments(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	envs, err := h.svc.EveAPI.GetEnvironments(ctx)
	if err != nil {
		if resourceNotFoundError(err) {
			h.svc.ChatService.UserNotificationThread(ctx, "failed to get environments", cmd.Info().User, cmd.Info().Channel, *ts)
			return
		}
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	if envs == nil {
		h.svc.ChatService.UserNotificationThread(ctx, "no environments", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	h.svc.ChatService.ShowResultsMessageThread(ctx, eveapi.ChatMessage(envs), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h ShowHandler) showNamespaces(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, err := h.svc.EveAPI.GetNamespacesByEnvironment(ctx, cmd.Options()[params.EnvironmentName].(string))
	if err != nil {
		if resourceNotFoundError(err) {
			h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("failed to get namespaces in environment: %s", cmd.Options()[params.EnvironmentName].(string)), cmd.Info().User, cmd.Info().Channel, *ts)
			return
		}
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	if ns == nil {
		h.svc.ChatService.UserNotificationThread(ctx, "no namespaces", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	h.svc.ChatService.ShowResultsMessageThread(ctx, eveapi.ChatMessage(ns), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h ShowHandler) showJobs(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, err := resolveNamespace(ctx, h.svc.EveAPI, cmd)
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, "invalid environment namespace request", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	nsJobs, err := h.svc.EveAPI.GetNamespaceJobs(ctx, &ns)
	if err != nil {
		if resourceNotFoundError(err) {
			h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("no jobs found for namespace: %s", ns.Alias), cmd.Info().User, cmd.Info().Channel, *ts)
			return
		}
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	if len(nsJobs) == 0 {
		h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("no jobs found for namespace: %s", ns.Alias), cmd.Info().User, cmd.Info().Channel, *ts)
	}
	h.svc.ChatService.ShowResultsMessageThread(ctx, eveapi.ChatMessage(nsJobs), cmd.Info().User, cmd.Info().Channel, *ts)
}

func (h ShowHandler) showServices(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	nv, err := resolveNamespace(ctx, h.svc.EveAPI, cmd)
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, err.Error(), cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	svcs, err := h.svc.EveAPI.GetServicesByNamespace(ctx, nv.Name)
	if err != nil {
		if resourceNotFoundError(err) {
			h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("failed to get services in namespace: %s", nv.Name), cmd.Info().User, cmd.Info().Channel, *ts)
			return
		}
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, *ts, err)
		return
	}
	if svcs == nil {
		h.svc.ChatService.UserNotificationThread(ctx, "no services", cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}
	h.svc.ChatService.ShowResultsMessageThread(ctx, eveapi.ChatMessage(svcs), cmd.Info().User, cmd.Info().Channel, *ts)
}

func resourceNotFoundError(err error) bool {
	if e, ok := err.(errors.RestError); ok {
		if e.Code == 404 {
			return true
		}
	}
	return false
}

func (h ShowHandler) showMetadata(ctx context.Context, cmd commands.EvebotCommand, ts *string) {
	ns, svc := resolveServiceNamespace(ctx, h.svc.EveAPI, h.svc.ChatService, cmd, ts)
	if svc == nil || ns == nil {
		return
	}

	mdKey := metaDataServiceKey(svc.Name, ns.Name)

	metadata, err := h.svc.EveAPI.GetMetadata(ctx, mdKey)
	if err != nil {
		h.svc.ChatService.UserNotificationThread(ctx, fmt.Sprintf("no metadata found for: %s", mdKey), cmd.Info().User, cmd.Info().Channel, *ts)
		return
	}

	h.svc.ChatService.ShowResultsMessageThread(ctx, eveapi.ChatMessage(metadata), cmd.Info().User, cmd.Info().Channel, *ts)
}
