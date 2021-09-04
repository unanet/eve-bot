package executor

import (
	"context"
	"errors"

	"github.com/unanet/eve-bot/internal/botcommander/interfaces"
	"github.com/unanet/eve-bot/internal/service"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/commands/handlers"
)

// EvebotCommandExecutor is the data structure that implements the Executor
type EvebotCommandExecutor struct {
	svc               *service.Provider
	cmdHandlerFactory handlers.Factory
}

// New creates a new executor
func New(svc *service.Provider, handlerFactory handlers.Factory) interfaces.CommandExecutor {
	return &EvebotCommandExecutor{
		svc:               svc,
		cmdHandlerFactory: handlerFactory,
	}
}

// Execute satisfies the Executor.Execute interface
func (h *EvebotCommandExecutor) Execute(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	if cmdHandlerFunc := h.cmdHandlerFactory.Items()[cmd.Info().CommandName]; cmdHandlerFunc != nil {
		cmdHandlerFunc(h.svc).Handle(ctx, cmd, timestamp)
		return
	}
	h.svc.ChatService.ErrorNotificationThread(ctx, cmd.Info().User, cmd.Info().Channel, timestamp, errors.New("failed to execute command; invalid command handler"))
}
