package executor

import (
	"context"
	"fmt"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands/handlers"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// Executor interface takes an EvebotCommand and Executes a matching handler
type Executor interface {
	Execute(ctx context.Context, cmd commands.EvebotCommand, timestamp string)
}

// EvebotCommandExecutor is the data structure that implements the Executor
type EvebotCommandExecutor struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
	cmdFactory   handlers.Factory
}

// New creates a new executor
func New(eveAPIClient eveapi.Client, chatSVC chatservice.Provider, handlerFactor handlers.Factory) Executor {
	return &EvebotCommandExecutor{
		eveAPIClient: eveAPIClient,
		chatSvc:      chatSVC,
		cmdFactory:   handlerFactor,
	}
}

// Execute satisfies the Executor.Execute interface
func (h *EvebotCommandExecutor) Execute(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	if cmdHandlerFunc := h.cmdFactory.Items()[cmd.Info().CommandName]; cmdHandlerFunc != nil {
		cmdHandlerFunc(&h.eveAPIClient, &h.chatSvc).Handle(ctx, cmd, timestamp)
		return
	}
	h.invalidCommandHandlerErr(ctx, "nil handler", cmd.Info().Channel, timestamp)
	return
}

func cleanSlackMsg(msg string) string {
	return fmt.Sprintf("\n\n ```%s```\n\n", msg)
}

func (h *EvebotCommandExecutor) invalidCommandHandlerErr(ctx context.Context, msg, channel, ts string) {
	err := fmt.Errorf("invalid command handler: %s", msg)
	log.Logger.Error("invalid command handler", zap.Error(err))
	h.chatSvc.PostMessageThread(ctx, cleanSlackMsg(err.Error()), channel, ts)
}
