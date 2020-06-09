package botcommander

import (
	"context"
	"fmt"
	"reflect"

	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands/handlers"
	"gitlab.unanet.io/devops/eve/pkg/log"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

type Executor interface {
	Execute(ctx context.Context, cmd botcommands.EvebotCommand, timestamp string)
}

type EvebotCommandExecutor struct {
	eveAPIClient eveapi.Client
	chatSvc      chatservice.Provider
}

func NewExecutor(eveAPIClient eveapi.Client, chatSVC chatservice.Provider) Executor {
	return &EvebotCommandExecutor{
		eveAPIClient: eveAPIClient,
		chatSvc:      chatSVC,
	}
}

func (h *EvebotCommandExecutor) Execute(ctx context.Context, cmd botcommands.EvebotCommand, timestamp string) {
	log.Logger.Debug("command handler execute", zap.Any("cmd_type", reflect.TypeOf(cmd)))
	cmdHandlerFunc := handlers.CommandHandlerMap[cmd.Name()]
	if cmdHandlerFunc == nil {
		h.invalidCommandHandlerErr(ctx, "nil handler", cmd.Channel(), timestamp)
		return
	}
	if cmdHandlerFuncVal, ok := cmdHandlerFunc.(func(*eveapi.Client, *chatservice.Provider) handlers.CommandHandler); ok {
		cmdHandlerFuncVal(&h.eveAPIClient, &h.chatSvc).Handle(ctx, cmd, timestamp)
		return
	}
	h.invalidCommandHandlerErr(ctx, "failed command type cast", cmd.Channel(), timestamp)
}

func cleanSlackMsg(msg string) string {
	return fmt.Sprintf("\n\n ```%s```\n\n", msg)
}

func (h *EvebotCommandExecutor) invalidCommandHandlerErr(ctx context.Context, msg, channel, ts string) {
	err := fmt.Errorf("invalid command handler: %s", msg)
	log.Logger.Error("invalid command handler", zap.Error(err))
	h.chatSvc.PostMessageThread(ctx, cleanSlackMsg(err.Error()), channel, ts)
}
