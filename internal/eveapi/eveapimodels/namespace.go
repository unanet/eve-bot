package eveapimodels

import (
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

type Namespaces []eve.Namespace

func (e Namespaces) ToChatMessage() string {
	if e == nil || len(e) == 0 {
		return "no environments"
	}

	msg := ""

	for _, v := range e {
		msg += "`" + v.Alias + "` ( _" + v.RequestedVersion + "_ )" + "\n"
	}

	return msg
}
