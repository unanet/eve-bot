package handlers

import (
	"context"

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
	h.chatSvc.UserNotificationThread(ctx, cmd.APIOptions()["resource"].(string), cmd.User(), cmd.Channel(), timestamp)
	h.chatSvc.UserNotificationThread(ctx, cmd.APIOptions()[params.ServiceName].(string), cmd.User(), cmd.Channel(), timestamp)
	h.chatSvc.UserNotificationThread(ctx, cmd.APIOptions()[params.NamespaceName].(string), cmd.User(), cmd.Channel(), timestamp)
	h.chatSvc.UserNotificationThread(ctx, cmd.APIOptions()[params.EnvironmentName].(string), cmd.User(), cmd.Channel(), timestamp)

	metadata := cmd.APIOptions()[params.MetadataName].(params.MetadataMap)

	h.chatSvc.UserNotificationThread(ctx, metadata.ToString(), cmd.User(), cmd.Channel(), timestamp)

}
