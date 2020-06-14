package resources

import "strings"

type Resource interface {
	Name() string
	Description() string
	Value() string
}

type Resources []Resource

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

var validResources = []Resource{DefaultEnvironment(), DefaultNamespace(), DefaultService()}

func IsValid(res string) bool {
	for _, v := range validResources {
		if v.Name() == strings.ToLower(res) {
			return true
		}
	}
	return false
}
