package botcommandhandlers

import (
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botcommands"
	"gitlab.unanet.io/devops/eve-bot/internal/chatservice"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

type DeployHandler struct {
	eveAPIClient *eveapi.Client
	chatSvc      *chatservice.Provider
}

func NewDeployHandler(*eveapi.Client, *chatservice.Provider) CommandHandler {
	return DeployHandler{}
}

func (h DeployHandler) Handle(cmd botcommands.EvebotCommand, timestamp string) {

}
