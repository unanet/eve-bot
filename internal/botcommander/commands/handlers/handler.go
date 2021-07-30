package handlers

import (
	"context"
	"errors"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
)

var errInvalidAPIResp = errors.New("invalid api response")

// CommandHandler is the interface that Handles EvebotCommands
type CommandHandler interface {
	Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string)
}
