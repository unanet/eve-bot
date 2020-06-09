package handlers

import (
	"context"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/log"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

type DeployHandler struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewDeployHandler(eveAPIClient *eveapi.Client, chatSvc *chatservice.Provider) CommandHandler {
	return DeployHandler{
		eveAPIClient: *eveAPIClient,
		chatSvc:      *chatSvc,
	}
}

func (h DeployHandler) Handle(ctx context.Context, cmd botcommands.EvebotCommand, timestamp string) {
	log.Logger.Debug("Are we there yet...")
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
