package commander

import "strings"

/*
	ARGUMENT: Services
*/

type EvebotArgService struct {
	Name, Version string
}

type EvebotArgServices []EvebotArgService

func (ebas EvebotArgServices) Name() string {
	return "services"
}

func NewParsedArgServices(input []string) EvebotArgServices {
	var svcs = EvebotArgServices{}
	for _, v := range input {
		kv := strings.Split(v, ":")
		if len(kv) > 1 {
			svcs = append(svcs, EvebotArgService{
				Name:    kv[0],
				Version: kv[1],
			})
		} else {
			svcs = append(svcs, EvebotArgService{
				Name:    kv[0],
				Version: "",
			})
		}
	}
	return svcs
}

func (ebas EvebotArgServices) Description() string {
	return "comma separated list of services with name:version syntax"
}

func NewServicesArg() EvebotArgServices {
	return EvebotArgServices{}
}
