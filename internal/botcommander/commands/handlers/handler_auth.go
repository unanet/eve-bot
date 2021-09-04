package handlers

import (
	"context"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/service"
)

// AuthHandler is the handler for the AuthCmd
type AuthHandler struct {
	svc *service.Provider
}

// NewAuthHandler creates a AuthHandler
func NewAuthHandler(svc *service.Provider) CommandHandler {
	return AuthHandler{svc: svc}
}

// Handle handles the AuthCmd
func (h AuthHandler) Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	chatUser, err := h.svc.ChatService.GetUser(ctx, cmd.Info().User)
	if err != nil {
		h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, err)
		return
	}
	h.svc.ChatService.PostPrivateMessage(ctx, h.svc.AuthCodeURL(chatUser.FullyQualifiedName()), cmd.Info().User)
}
