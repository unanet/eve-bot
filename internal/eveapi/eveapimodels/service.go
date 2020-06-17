package eveapimodels

import (
	"strconv"

	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Service struct {
	ID              int       `json:"id"`
	NamespaceID     int       `json:"namespace_id"`
	NamespaceName   string    `json:"namespace_name"`
	ArtifactID      int       `json:"artifact_id"`
	OverrideVersion string    `json:"override_version"`
	DeployedVersion string    `json:"deployed_version"`
	Metadata        json.Text `json:"metadata"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	Name            string    `json:"name"`
	StickySessions  bool      `json:"sticky_sessions"`
	Count           int       `json:"count"`
}

type Services []Service

func (e Services) ToChatMessage() string {
	if e == nil || len(e) == 0 {
		return "no services"
	}

	msg := ""

	for _, v := range e {
		msg += "*Name:* " + "`" + v.Name + "`" + "		" + "*Version:* " + "_" + v.DeployedVersion + "_" + "\n" +
			"*Sticky:* " + strconv.FormatBool(v.StickySessions) + "		" + "*Count:* " + strconv.Itoa(v.Count) + "\n\n"
	}

	return msg
}
