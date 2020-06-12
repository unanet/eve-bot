package handlers

import (
	"context"
	"errors"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
)

var (
	errInvalidApiResp = errors.New("invalid api response")
)

type CommandHandler interface {
	Handle(ctx context.Context, cmd commands.EvebotCommand, timestamp string)
}
