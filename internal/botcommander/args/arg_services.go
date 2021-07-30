package args

import (
	"strings"

	"github.com/unanet/eve/pkg/eve"
)

/*
	ARGUMENT: Services
*/

const (
	// ServicesName is the Services argument name
	ServicesName = "services"
	// ServicesDescription us the Services argument description
	ServicesDescription = "comma separated list of services with name:version syntax (version is optional)"
)

// Service is the Service argument type
type Service struct {
	Name, Version string
}

// Services is a slice of Service arguments
type Services []Service

// Name is the name os the Services argument
func (svcs Services) Name() string {
	return ServicesName
}

// Description is the description of the Services argument
func (svcs Services) Description() string {
	return ServicesDescription
}

// Value is the value of the Services argument
func (svcs Services) Value() interface{} {
	var artifactDefs eve.ArtifactDefinitions

	for _, v := range svcs {
		artifactDefs = append(artifactDefs, &eve.ArtifactDefinition{
			Name:             v.Name,
			RequestedVersion: v.Version,
		})
	}

	return artifactDefs
}

// DefaultServicesArg is the default Services argument
func DefaultServicesArg() Services {
	return Services{}
}

// NewServicesArg is the instantiation method that creates a new Services argument
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
