package botcommands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

type EvebotDeployCommand struct {
	baseCommand
}

func NewEvebotDeployCommand() EvebotDeployCommand {
	return EvebotDeployCommand{baseCommand{
		name:    "deploy",
		summary: "Deploy command is used to deploy services to a specific namespace and environment",
		usage: bothelp.HelpUsage{
			"deploy {{ namespace }} in {{ environment }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		optionalArgs:   botargs.Args{botargs.NewDryrunArg(), botargs.NewForceArg(), botargs.NewServicesArg()},
		additionalArgs: botargs.Args{},
		examples: bothelp.HelpExamples{
			"deploy current in qa",
			"deploy current in qa services=infocus-cloud-client:2020.1 dryrun=true",
			"deploy current in qa services=infocus-cloud-client:2020.1,infocus-proxy:2020.1 dryrun=true force=true",
			"deploy current in qa services=infocus-cloud-client,infocus-proxy",
		},
		asyncRequired: true,
	}}
}

func (cmd EvebotDeployCommand) AsyncRequired() bool {
	return cmd.asyncRequired
}

func (cmd EvebotDeployCommand) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd EvebotDeployCommand) Name() string {
	return cmd.name
}

func (cmd EvebotDeployCommand) Help() *bothelp.Help {
	return bothelp.NewEvebotCommandHelp(
		bothelp.EvebotCommandHelpSummaryOpt(cmd.summary.String()),
		bothelp.EvebotCommandHelpUsageOpt(cmd.usage.String()),
		bothelp.EvebotCommandHelpArgsOpt(cmd.optionalArgs.String()),
		bothelp.EvebotCommandHelpExamplesOpt(cmd.examples.String()),
	)
}

func (cmd EvebotDeployCommand) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd EvebotDeployCommand) AdditionalArgs() (botargs.Args, error) {
	// if we've already calculated the args, use them
	if len(cmd.additionalArgs) > 0 {
		return cmd.additionalArgs, nil
	}

	// haven't calculated the args and no need since they weren't supplied
	if len(cmd.input) <= 3 {
		return botargs.Args{}, nil
	}

	// let's calculate the args based on the input command
	var additionalArgs botargs.Args

	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if additionalArg := botargs.ResolveArgumentKV(argKV); additionalArg != nil {
				additionalArgs = append(additionalArgs, additionalArg)
			} else {
				return botargs.Args{}, fmt.Errorf("invalid additional arg: %v", argKV)
			}
		}
	}

	// set the calculated args so we don't have to calculate them again
	cmd.additionalArgs = additionalArgs
	return additionalArgs, nil
}
