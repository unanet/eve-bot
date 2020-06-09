package botcommander

import (
	"context"
	"errors"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommandhandlers"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

var (
	errInvalidRequestObj = errors.New("invalid request object")
	errInvalidApiResp    = errors.New("invalid api response")
)

type Handler interface {
	Handle(ctx context.Context, cmd botcommands.EvebotCommand, timestamp string)
	Resolve(ctx context.Context, cmd botcommands.EvebotCommand, timestamp string)
}

type EvebotCommandHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewHandler(eveAPIClient eveapi.Client, chatSVC chatservice.Provider) Handler {
	return &EvebotCommandHandler{
		eveAPIClient: eveAPIClient,
		chatSvc:      chatSVC,
	}
}

func (h *EvebotCommandHandler) Resolve(ctx context.Context, cmd botcommands.EvebotCommand, timestamp string) {
	switch cmd.(type) {
	case botcommands.DeployCmd:
		botcommandhandlers.NewDeployHandler(&h.eveAPIClient, &h.chatSvc).Handle(cmd, timestamp)
	case botcommands.MigrateCmd:
	}
}

func (h *EvebotCommandHandler) Handle(ctx context.Context, cmd botcommands.EvebotCommand, timestamp string) {
	chatUser, err := h.chatSvc.GetUser(ctx, cmd.User())
	if err != nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, err)
		return
	}
	var reqObj interface{}
	if reqObj = cmd.EveReqObj(chatUser.Name); reqObj == nil {
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, errInvalidRequestObj)
		return
	}
	switch reqObj.(type) {
	case eveapi.DeploymentPlanOptions:
		resp, err := h.eveAPIClient.Deploy(ctx, reqObj.(eveapi.DeploymentPlanOptions), cmd.User(), cmd.Channel(), timestamp)
		if err != nil && len(err.Error()) > 0 {
			h.chatSvc.DeploymentNotificationThread(ctx, err.Error(), cmd.User(), cmd.Channel(), timestamp)
			return
		}
		if resp == nil {
			h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, errInvalidApiResp)
			return
		}
		if len(resp.Messages) > 0 {
			h.chatSvc.UserNotificationThread(ctx, strings.Join(resp.Messages, ","), cmd.User(), cmd.Channel(), timestamp)
			return
		}
		return
	default:
		h.chatSvc.ErrorNotificationThread(ctx, cmd.User(), cmd.Channel(), timestamp, errInvalidRequestObj)
		return
	}
}
