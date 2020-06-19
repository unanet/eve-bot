package eveapimodels

import (
	"encoding/json"

	"gitlab.unanet.io/devops/eve/pkg/eve"
)

type EveService eve.Service

type Services []eve.Service

func (e Services) ToChatMessage() string {
	if e == nil || len(e) == 0 {
		return "no services"
	}

	msg := ""

	for _, v := range e {
		msg += "`" + v.Name + "` - _" + v.DeployedVersion + "_ ( *" + v.ArtifactName + "* )" + "\n"
	}

	return msg
}

func (s EveService) MetadataToChatMessage() string {
	if s.ID == 0 || len(s.Metadata) == 0 {
		return "no metadata"
	}

	jsonB, err := json.MarshalIndent(s.Metadata, "", "    ")
	if err != nil {
		return "invalid json metadata"
	}

	return "```" + string(jsonB) + "```"
}
