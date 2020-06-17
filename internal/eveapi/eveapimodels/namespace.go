package eveapimodels

import (
	"database/sql"

	"gitlab.unanet.io/devops/eve/pkg/json"
)

type Namespace struct {
	ID                 int          `db:"id"`
	Name               string       `db:"name"`
	Alias              string       `db:"alias"`
	EnvironmentID      int          `db:"environment_id"`
	EnvironmentName    string       `db:"environment_name"`
	RequestedVersion   string       `db:"requested_version"`
	ExplicitDeployOnly bool         `db:"explicit_deploy_only"`
	ClusterID          int          `db:"cluster_id"`
	Metadata           json.Text    `db:"metadata"`
	CreatedAt          sql.NullTime `db:"created_at"`
	UpdatedAt          sql.NullTime `db:"updated_at"`
}

type Namespaces []Namespace

func (e Namespaces) ToChatMessage() string {
	if e == nil || len(e) == 0 {
		return "no environments"
	}

	msg := ""

	for _, v := range e {
		msg += v.Alias + " (" + v.RequestedVersion + ")" + "\n\n"
	}

	return msg
}
