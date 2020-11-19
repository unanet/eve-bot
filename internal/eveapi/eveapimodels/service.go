package eveapimodels

import (
	"encoding/json"

	"gitlab.unanet.io/devops/eve/pkg/eve"
)

// EveService data structure
type EveService eve.Service

// Services data structure
type Services []eve.Service

// ToChatMessage converts the Services data structure to a Formatted Chat Message
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

// MetadataToChatMessage converts the Service.Metadata to chat message
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

type MetaData struct {
	Input eve.Metadata
}

func (d MetaData) ToChatMessage() string {
	if d.Input.ID == 0 || len(d.Input.Value) == 0 {
		return "no metadata"
	}

	jsonB, err := json.MarshalIndent(d.Input.Value, "", "	")
	if err != nil {
		return "invalid json metadata"
	}
	return "```" + string(jsonB) + "```"
}
