package commander

import (
	"fmt"
	"strconv"
	"strings"
)

type EvebotDeployCommand struct {
	name string
}

func NewEvebotDeployCommand() *EvebotDeployCommand {
	return &EvebotDeployCommand{
		name: "deploy",
	}
}

func (edc *EvebotDeployCommand) Name() string {
	return edc.name
}

func (edc *EvebotDeployCommand) Examples() EvebotCommandExamples {
	return EvebotCommandExamples{
		"- deploy {{ namespace }} in {{ environment }}",
		"- deploy {{ namespace }} in {{ environment }} services={{ artifact_name:artifact_version }}",
		"- deploy {{ namespace }} in {{ environment }} services={{ artifact_name:artifact_version }} dryrun={{ true }}",
		"- deploy {{ namespace }} in {{ environment }} services={{ artifact_name:artifact_version }} dryrun={{ true }} force={{ true }}",
		"\n",
		"`Examples:`",
		"- deploy current in qa",
		"- deploy current in qa services=infocus-cloud-client:2020.1 dryrun=true",
		"- deploy current in qa services=infocus-cloud-client:2020.1,infocus-proxy:2020.1 dryrun=true redeploy=true",
	}
}

func (edc *EvebotDeployCommand) IsHelpRequest(input []string) bool {
	if input[1] == "help" {
		return true
	}
	return false
}

func (edc *EvebotDeployCommand) IsValidCommand(input []string) bool {
	if len(input) <= 3 || input[0] != edc.Name() {
		return false
	}
	return true
}

func (edc *EvebotDeployCommand) AdditionalArgs(input []string) (EvebotArgs, error) {
	if len(input) <= 3 {
		return EvebotArgs{}, nil
	}

	var additionalArgs EvebotArgs

	for _, s := range input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if additionalArg := edc.ResolveAdditionalArg(argKV); additionalArg != nil {
				additionalArgs = append(additionalArgs, additionalArg)
			} else {
				return EvebotArgs{}, fmt.Errorf("invalid additional arg: %v", argKV)
			}
		}
	}

	return additionalArgs, nil
}

func (edc *EvebotDeployCommand) ResolveAdditionalArg(argKV []string) EvebotArg {
	switch strings.ToLower(argKV[0]) {
	case "dryrun":
		b, err := strconv.ParseBool(strings.ToLower(argKV[1]))
		if err != nil {
			return nil
		} else {
			return EvebotArgDryrun(b)
		}
	case "force":
		b, err := strconv.ParseBool(strings.ToLower(argKV[1]))
		if err != nil {
			return nil
		} else {
			return EvebotArgForce(b)
		}
	case "services":
		return EvebotArgServices(strings.Split(argKV[1], ","))
	default:
		return nil
	}

}
