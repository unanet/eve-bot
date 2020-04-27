package commander

import "strings"

/*
	ARGUMENT: Services
*/

type ArgService struct {
	Name, Version string
}

type ArgServices []ArgService

func (ebas ArgServices) Name() string {
	return "services"
}

func NewParsedArgServices(input []string) ArgServices {
	var svcs = ArgServices{}
	for _, v := range input {
		kv := strings.Split(v, ":")
		if len(kv) > 1 {
			svcs = append(svcs, ArgService{
				Name:    kv[0],
				Version: kv[1],
			})
		} else {
			svcs = append(svcs, ArgService{
				Name:    kv[0],
				Version: "",
			})
		}
	}
	return svcs
}

func (ebas ArgServices) Description() string {
	return "comma separated list of services with name:version and version is optional"
}

func NewServicesArg() ArgServices {
	return ArgServices{}
}
