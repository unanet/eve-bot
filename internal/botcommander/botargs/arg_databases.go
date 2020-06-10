package botargs

import (
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

/*
	ARGUMENT: Databases
*/

const (
	DatabasesName = "databases"
)

type Database struct {
	Name, Version string
}

type Databases []Database

func (a Databases) Name() string {
	return DatabasesName
}

func (a Databases) Description() string {
	return "comma separated list of databases with name:version syntax (version is optional)"
}

func (a Databases) Value() interface{} {
	var artifactDefs eveapi.ArtifactDefinitions

	for _, v := range a {
		artifactDefs = append(artifactDefs, &eveapi.ArtifactDefinition{
			Name:             v.Name,
			RequestedVersion: v.Version,
		})
	}

	return artifactDefs
}

func DefaultDatabasesArg() Databases {
	return Databases{}
}

func NewDatabasesArg(input []string) Databases {
	var dbs = Databases{}
	for _, v := range input {
		kv := strings.Split(v, ":")
		if len(kv) > 1 {
			dbs = append(dbs, Database{
				Name:    kv[0],
				Version: kv[1],
			})
		} else {
			dbs = append(dbs, Database{
				Name:    kv[0],
				Version: "",
			})
		}
	}
	return dbs
}
