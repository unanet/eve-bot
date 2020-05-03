package botargs

import (
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

/*
	ARGUMENT: Services
*/

const (
	ServicesName = "services"
)

type Service struct {
	Name, Version string
}

type Services []Service

func (ebas Services) Name() string {
	return ServicesName
}

func (ebas Services) Description() string {
	return "comma separated list of services with name:version syntax (version is optional)"
}

func (ebas Services) Value() interface{} {
	var artifactDefs []eveapi.ArtifactDefinition

	for _, v := range ebas {
		artifactDefs = append(artifactDefs, eveapi.ArtifactDefinition{
			Name:             v.Name,
			RequestedVersion: v.Version,
		})
	}

	return artifactDefs
}

func DefaultServicesArg() Services {
	return Services{}
}

func NewServicesArg(input []string) Services {
	var svcs = Services{}
	for _, v := range input {
		kv := strings.Split(v, ":")
		if len(kv) > 1 {
			svcs = append(svcs, Service{
				Name:    kv[0],
				Version: kv[1],
			})
		} else {
			svcs = append(svcs, Service{
				Name:    kv[0],
				Version: "",
			})
		}
	}
	return svcs
}
