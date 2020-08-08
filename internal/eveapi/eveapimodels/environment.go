package eveapimodels

import (
	"gitlab.unanet.io/devops/eve/pkg/eve"
)

// Environments data structure
type Environments []eve.Environment

// ToChatMessage converts environments to a chat message
func (e Environments) ToChatMessage() string {
	if e == nil || len(e) == 0 {
		return "no environments"
	}

	msg := ""

	for _, v := range e {
		msg += "*Name:* " + "`" + v.Name + "`" + "\n" + "*Description:* " + "_" + v.Description + "_" + "\n\n"
	}

	return msg
}
