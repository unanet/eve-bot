package eveapimodels

import "gitlab.unanet.io/devops/eve/pkg/json"

type Environment struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Alias       string    `db:"alias"`
	Description string    `db:"description"`
	Metadata    json.Text `db:"metadata"`
}

type Environments []Environment

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
