package resources

import "strings"

// Resource interface
type Resource interface {
	Name() string
	Description() string
	Value() string
}

// Resources slice of resource
type Resources []Resource

// String satisfies interface and converts slice of resources to a string
func (p Resources) String() string {
	var msg string
	for _, v := range p {
		msg = msg + v.Name() + " - " + v.Description() + "\n"
	}
	return msg
}

type baseResource struct {
	name        string
	description string
	value       string
}

// FullResourceMap is just a map of resources that are available
// This map should never be written to, just read for Validation
var FullResourceMap = map[string]bool{
	strings.ToLower(EnvironmentName): true,
	strings.ToLower(NamespaceName):   true,
	strings.ToLower(ServiceName):     true,
	strings.ToLower(MetadataName):    true,
	strings.ToLower(JobName):         true,
	"jobs":                           true, // Job vs Jobs TODO: Clean this up
	strings.ToLower(VersionName):     true,
}

// ValidResMutations are just a map of resources that can be mutated by the bot (user)
// This map should never be written to, just read for Validation
var ValidResourcesMutations = map[string]bool{
	strings.ToLower(MetadataName): true,
	strings.ToLower(VersionName):  true,
}
