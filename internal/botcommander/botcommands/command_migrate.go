package botcommands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

type EvebotMigrateCommand struct {
	baseCommand
}

func NewEvebotMigrateCommand() EvebotDeployCommand {
	return EvebotDeployCommand{baseCommand{
		name:    "migrate",
		summary: "Migrate command is used to migrate databases by namespace and environment",
		usage: bothelp.HelpUsage{
			"migrate {{ namespace }} in {{ environment }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ service_name:service_version }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		optionalArgs:   botargs.Args{botargs.NewDryrunArg(), botargs.NewForceArg(), botargs.NewDatabasesArg()},
		additionalArgs: botargs.Args{},
		examples: bothelp.HelpExamples{
			"migrate current in qa",
			"migrate current in qa databases=infocus dryrun=true",
			"migrate current in qa databases=infocus dryrun=true force=true",
			"migrate current in qa databases=infocus,cloud-support",
		},
		asyncRequired: true,
	}}
}

func (cmd EvebotMigrateCommand) AsyncRequired() bool {
	return cmd.asyncRequired
}

func (cmd EvebotMigrateCommand) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd EvebotMigrateCommand) Name() string {
	return cmd.name
}

func (cmd EvebotMigrateCommand) Help() *bothelp.Help {
	return bothelp.NewEvebotCommandHelp(
		bothelp.EvebotCommandHelpSummaryOpt(cmd.summary.String()),
		bothelp.EvebotCommandHelpUsageOpt(cmd.usage.String()),
		bothelp.EvebotCommandHelpArgsOpt(cmd.optionalArgs.String()),
		bothelp.EvebotCommandHelpExamplesOpt(cmd.examples.String()),
	)
}

func (cmd EvebotMigrateCommand) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd EvebotMigrateCommand) AdditionalArgs() (botargs.Args, error) {
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
