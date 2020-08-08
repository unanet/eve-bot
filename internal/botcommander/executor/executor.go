package executor

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands/handlers"
	"gitlab.unanet.io/devops/eve/pkg/log"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/commands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

// Executor interface takes an EvebotCommand and Executes a matching handler
type Executor interface {
	Execute(ctx context.Context, cmd commands.EvebotCommand, timestamp string)
}

// EvebotCommandExecutor is the data structure that implements the Executor
type EvebotCommandExecutor struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

// NewExecutor creator a new executor
func NewExecutor(eveAPIClient eveapi.Client, chatSVC chatservice.Provider) Executor {
	return &EvebotCommandExecutor{
		eveAPIClient: eveAPIClient,
		chatSvc:      chatSVC,
	}
}

// Execute satisfies the Executor.Execute interface
func (h *EvebotCommandExecutor) Execute(ctx context.Context, cmd commands.EvebotCommand, timestamp string) {
	cmdHandlerFunc := handlers.CommandHandlerMap[cmd.Info().CommandName]
	if cmdHandlerFunc == nil {
		h.invalidCommandHandlerErr(ctx, "nil handler", cmd.Info().Channel, timestamp)
		return
	}
	if cmdHandlerFuncVal, ok := cmdHandlerFunc.(func(*eveapi.Client, *chatservice.Provider) handlers.CommandHandler); ok {
		cmdHandlerFuncVal(&h.eveAPIClient, &h.chatSvc).Handle(ctx, cmd, timestamp)
		return
	}
	h.invalidCommandHandlerErr(ctx, "failed command type cast", cmd.Info().Channel, timestamp)
}

func cleanSlackMsg(msg string) string {
	return fmt.Sprintf("\n\n ```%s```\n\n", msg)
}

func (h *EvebotCommandExecutor) invalidCommandHandlerErr(ctx context.Context, msg, channel, ts string) {
	err := fmt.Errorf("invalid command handler: %s", msg)
	log.Logger.Error("invalid command handler", zap.Error(err))
	h.chatSvc.PostMessageThread(ctx, cleanSlackMsg(err.Error()), channel, ts)
}
