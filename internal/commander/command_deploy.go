package commander

import (
	"fmt"
	"strconv"
	"strings"
)

type EvebotDeployCommand struct {
	name         string
	input        []string
	summary      EvebotCommandSummary
	usage        EvebotCommandUsage
	optionalArgs EvebotCommandArgs
	examples     EvebotCommandExamples
}

func NewEvebotDeployCommand() EvebotDeployCommand {
	return EvebotDeployCommand{
		name:    "deploy",
		summary: "Deploy command is used to deploy services to a specific namespace and environment",
		usage: EvebotCommandUsage{
			"deploy {{ namespace }} in {{ environment }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		optionalArgs: EvebotCommandArgs{NewDryrunArg(), NewForceArg(), NewServicesArg()},
		examples: EvebotCommandExamples{
			"deploy current in qa",
			"deploy current in qa services=infocus-cloud-client:2020.1 dryrun=true",
			"deploy current in qa services=infocus-cloud-client:2020.1,infocus-proxy:2020.1 dryrun=true force=true",
		},
	}
}

func (edc EvebotDeployCommand) Initialize(input []string) EvebotCommand {
	edc.input = input
	return edc
}

func (edc EvebotDeployCommand) Name() string {
	return edc.name
}

func (edc EvebotDeployCommand) Help() *EvebotCommandHelp {
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

func (edc EvebotDeployCommand) AdditionalArgs() (EvebotCommandArgs, error) {
	if len(edc.input) <= 3 {
		return EvebotCommandArgs{}, nil
	}

	var additionalArgs EvebotCommandArgs

	for _, s := range edc.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if additionalArg := edc.resolveAdditionalArg(argKV); additionalArg != nil {
				additionalArgs = append(additionalArgs, additionalArg)
			} else {
				return EvebotCommandArgs{}, fmt.Errorf("invalid additional arg: %v", argKV)
			}
		}
	}

	return additionalArgs, nil
}

func (edc EvebotDeployCommand) resolveAdditionalArg(argKV []string) EvebotCommandArg {
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
		return NewParsedArgServices(strings.Split(argKV[1], ","))
	default:
		return nil
	}

}
