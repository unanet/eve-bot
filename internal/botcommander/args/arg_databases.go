package args

import (
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi/eveapimodels"
)

/*
	ARGUMENT: Databases
*/

const (
	// DatabasesName argument name
	DatabasesName = "databases"
	// DatabaseDescription describes the database argument and it's purpose
	DatabaseDescription = "comma separated list of databases with name:version syntax (version is optional)"
)

// Database argument structure
type Database struct {
	Name, Version string
}

// Databases slice of database structure
type Databases []Database

// Name of the Databases argument
func (dbs Databases) Name() string {
	return DatabasesName
}

// Description of the Databases argument
func (dbs Databases) Description() string {
	return DatabaseDescription
}

// Value of the databases argument
func (dbs Databases) Value() interface{} {
	var artifactDefs eveapimodels.ArtifactDefinitions

	for _, v := range dbs {
		artifactDefs = append(artifactDefs, &eveapimodels.ArtifactDefinition{
			Name:             v.Name,
			RequestedVersion: v.Version,
		})
	}

	return artifactDefs
}

// DefaultDatabasesArg is the default databases argument
func DefaultDatabasesArg() Databases {
	return Databases{}
}

// NewDatabasesArg creates a new databases argument
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
