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

// IsValid validates id the supplied resource is valid
func IsValid(res string) bool {
	for _, v := range []Resource{DefaultEnvironment(), DefaultNamespace(), DefaultService(), DefaultMetadata()} {
		if v.Name() == strings.ToLower(res) {
			return true
		}
	}
	return false
}

// IsValidSet validates that the supplied resource can be set
func IsValidSet(res string) bool {
	for _, v := range []Resource{DefaultMetadata(), DefaultVersion()} {
		if v.Name() == strings.ToLower(res) {
			return true
		}
	}
	return false
}

// IsValidDelete validates that the supplied resource can be deleted
func IsValidDelete(res string) bool {
	for _, v := range []Resource{DefaultMetadata(), DefaultVersion()} {
		if v.Name() == strings.ToLower(res) {
			return true
		}
	}
	return false
}
