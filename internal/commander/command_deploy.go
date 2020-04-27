package commander

import (
	"fmt"
	"strconv"
	"strings"
)

type EvebotDeployCommand struct {
	input         []string
	name          string
	asyncRequired bool
	summary       HelpSummary
	usage         HelpUsage
	optionalArgs  Args
	examples      UserHelpExamples
}

func NewEvebotDeployCommand() EvebotDeployCommand {
	return EvebotDeployCommand{
		name:    "deploy",
		summary: "Deploy command is used to deploy services to a specific namespace and environment",
		usage: HelpUsage{
			"deploy {{ namespace }} in {{ environment }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		optionalArgs: Args{NewDryrunArg(), NewForceArg(), NewServicesArg()},
		examples: UserHelpExamples{
			"deploy current in qa",
			"deploy current in qa services=infocus-cloud-client:2020.1 dryrun=true",
			"deploy current in qa services=infocus-cloud-client:2020.1,infocus-proxy:2020.1 dryrun=true force=true",
			"deploy current in qa services=infocus-cloud-client,infocus-proxy",
		},
		asyncRequired: true,
	}
}

func (edc EvebotDeployCommand) AsyncRequired() bool {
	return edc.asyncRequired
}

func (edc EvebotDeployCommand) Initialize(input []string) EvebotCommand {
	edc.input = input
	return edc
}

func (edc EvebotDeployCommand) Name() string {
	return edc.name
}

func (edc EvebotDeployCommand) Help() *Help {
	return NewEvebotCommandHelp(
		EvebotCommandHelpSummaryOpt(edc.summary.String()),
		EvebotCommandHelpUsageOpt(edc.usage.String()),
		EvebotCommandHelpArgsOpt(edc.optionalArgs.String()),
		EvebotCommandHelpExamplesOpt(edc.examples.String()),
	)
}

func (edc EvebotDeployCommand) IsHelpRequest() bool {
	if edc.input[0] == "help" || edc.input[len(edc.input)-1] == "help" || (len(edc.input) == 1 && edc.input[0] == "deploy") {
		return true
	}
	return false
}

func (edc EvebotDeployCommand) AdditionalArgs() (Args, error) {
	if len(edc.input) <= 3 {
		return Args{}, nil
	}

	var additionalArgs Args

	for _, s := range edc.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if additionalArg := edc.resolveAdditionalArg(argKV); additionalArg != nil {
				additionalArgs = append(additionalArgs, additionalArg)
			} else {
				return Args{}, fmt.Errorf("invalid additional arg: %v", argKV)
			}
		}
	}

	return additionalArgs, nil
}

func (edc EvebotDeployCommand) resolveAdditionalArg(argKV []string) Arg {
	switch strings.ToLower(argKV[0]) {
	case "dryrun":
		b, err := strconv.ParseBool(strings.ToLower(argKV[1]))
		if err != nil {
			return nil
		} else {
			return ArgDryrun(b)
		}
	case "force":
		b, err := strconv.ParseBool(strings.ToLower(argKV[1]))
		if err != nil {
			return nil
		} else {
			return ArgForce(b)
		}
	case "services":
		return NewParsedArgServices(strings.Split(argKV[1], ","))
	default:
		return nil
	}

}
