package eveapimodels

import (
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

type Services []eve.Service

func (e Services) ToChatMessage() string {
	if e == nil || len(e) == 0 {
		return "no services"
	}

	msg := ""

	for _, v := range e {
		msg += "*" + v.Name + "* - _" + v.DeployedVersion + "_ ( _" + v.ArtifactName + "_ )" + "\n"
	}

	return msg
}
