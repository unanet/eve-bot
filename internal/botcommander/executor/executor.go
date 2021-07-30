package executor

import (
	"context"
	"errors"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/commands/handlers"
)

// EvebotCommandExecutor is the data structure that implements the Executor
type EvebotCommandExecutor struct {
	eveAPIClient interfaces.EveAPI
	chatSvc      interfaces.ChatProvider
	cmdFactory   handlers.Factory
}

// New creates a new executor
func New(eveAPIClient interfaces.EveAPI, chatSVC interfaces.ChatProvider, handlerFactor handlers.Factory) interfaces.CommandExecutor {
	return &EvebotCommandExecutor{
		eveAPIClient: eveAPIClient,
		chatSvc:      chatSVC,
		cmdFactory:   handlerFactor,
	}
}

// Execute satisfies the Executor.Execute interface
func (h *EvebotCommandExecutor) Execute(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	if cmdHandlerFunc := h.cmdFactory.Items()[cmd.Info().CommandName]; cmdHandlerFunc != nil {
		cmdHandlerFunc(h.eveAPIClient, h.chatSvc).Handle(ctx, cmd, timestamp)
		return
	}
	h.chatSvc.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, errors.New("failed to execute command; invalid command handler"))
	return
}
