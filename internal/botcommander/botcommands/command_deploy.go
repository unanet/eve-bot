package botcommands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

type DeployCmd struct {
	baseCommand
}

func NewEvebotDeployCommand() DeployCmd {
	return DeployCmd{baseCommand{
		name:    "deploy",
		summary: "The `deploy` command is used to deploy services to a specific *namespace* and *environment*",
		usage: bothelp.Usage{
			"deploy {{ namespace }} in {{ environment }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		optionalArgs:   botargs.Args{botargs.DefaultDryrunArg(), botargs.DefaultForceArg(), botargs.DefaultServicesArg()},
		additionalArgs: botargs.Args{},
		examples: bothelp.Examples{
			"deploy current in qa",
			"deploy current in qa services=infocus-cloud-client:2020.1 dryrun=true",
			"deploy current in qa services=infocus-cloud-client:2020.2.232,infocus-proxy:2020.2.199 dryrun=true force=true",
			"deploy current in qa services=infocus-cloud-client,infocus-proxy",
		},
		asyncRequired: true,
	}}
}

func (cmd DeployCmd) IsValid() bool {
	if cmd.input == nil || len(cmd.input) == 0 || len(cmd.input) <= 3 {
		return false
	}
	return true
}

func (cmd DeployCmd) AsyncRequired() bool {
	return cmd.asyncRequired
}

func (cmd DeployCmd) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd DeployCmd) Name() string {
	return cmd.name
}

func (cmd DeployCmd) Help() *bothelp.Help {
	return bothelp.New(
		bothelp.HeaderOpt(cmd.summary.String()),
		bothelp.UsageOpt(cmd.usage.String()),
		bothelp.ArgsOpt(cmd.optionalArgs.String()),
		bothelp.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd DeployCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd DeployCmd) AdditionalArgs() (botargs.Args, error) {
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
