package botcommandhandlers

import "gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"

type CommandHandler interface {
	Handle(cmd botcommands.EvebotCommand, timestamp string)
}
