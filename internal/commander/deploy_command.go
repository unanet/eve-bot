package commander

import (
	"fmt"
	"strconv"
	"strings"

	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

type EvebotDeployCommand struct {
	name         string
	summary      EvebotCommandSummary
	usage        EvebotCommandUsage
	optionalArgs EvebotCommandArgs
	examples     EvebotCommandExamples
	helpCmd      EvebotCommandHelp
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
		helpCmd: EvebotCommandHelp{},
	}
}

func (edc EvebotDeployCommand) Examples() EvebotCommandExamples {
	return edc.examples
}

func (edc EvebotDeployCommand) Name() string {
	return edc.name
}

func (edc EvebotDeployCommand) OptionalArgs() EvebotCommandArgs {
	return edc.optionalArgs
}

func (edc EvebotDeployCommand) Summary() EvebotCommandSummary {
	return edc.summary
}

func (edc EvebotDeployCommand) Usage() EvebotCommandUsage {
	return edc.usage
}

func (edc EvebotDeployCommand) Help() *EvebotCommandHelp {
	return NewEvebotCommandHelp(
		EvebotCommandHelpSummaryOpt(edc.summary.String()),
		EvebotCommandHelpUsageOpt(edc.usage.String()),
		EvebotCommandHelpArgsOpt(edc.optionalArgs.String()),
		EvebotCommandHelpExamplesOpt(edc.examples.String()),
	)
}

func (edc EvebotDeployCommand) IsHelpRequest(input []string) bool {
	log.Logger.Debug("input length", zap.Int("length", len(input)))
	if input[0] == "help" || input[len(input)-1] == "help" {
		return true
	}

	if len(input) == 1 && input[0] == "deploy" {
		return true
	}

	return false
}

func (edc EvebotDeployCommand) IsValidCommand(input []string) bool {
	if len(input) <= 3 || input[0] != edc.Name() {
		return false
	}
	return true
}

func (edc EvebotDeployCommand) AdditionalArgs(input []string) (EvebotCommandArgs, error) {
	if len(input) <= 3 {
		return EvebotCommandArgs{}, nil
	}

	var additionalArgs EvebotCommandArgs

	for _, s := range input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if additionalArg := edc.ResolveAdditionalArg(argKV); additionalArg != nil {
				additionalArgs = append(additionalArgs, additionalArg)
			} else {
				return EvebotCommandArgs{}, fmt.Errorf("invalid additional arg: %v", argKV)
			}
		}
	}

	return additionalArgs, nil
}

func (edc EvebotDeployCommand) ResolveAdditionalArg(argKV []string) EvebotArg {
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
