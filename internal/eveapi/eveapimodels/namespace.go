package eveapimodels

import (
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

// Namespaces data structure for Eve Namespaces
type Namespaces []eve.Namespace

// ToChatMessage converts the eve namespaces to a chat message
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
