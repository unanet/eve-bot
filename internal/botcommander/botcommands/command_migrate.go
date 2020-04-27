package botcommands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
)

type MigrateCmd struct {
	baseCommand
}

func NewEvebotMigrateCommand() DeployCmd {
	return DeployCmd{baseCommand{
		name:    "migrate",
		summary: "The `migrate` command is used to migrate databases by *namespace* and *environment*",
		usage: bothelp.Usage{
			"migrate {{ namespace }} in {{ environment }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ service_name:service_version }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		optionalArgs:   botargs.Args{botargs.DefaultDryrunArg(), botargs.DefaultForceArg(), botargs.DefaultDatabasesArg()},
		additionalArgs: botargs.Args{},
		examples: bothelp.Examples{
			"migrate current in qa",
			"migrate current in qa databases=infocus dryrun=true",
			"migrate current in qa databases=infocus dryrun=true force=true",
			"migrate current in qa databases=infocus,cloud-support",
		},
		asyncRequired: true,
	}}
}

func (cmd MigrateCmd) AsyncRequired() bool {
	return cmd.asyncRequired
}

func (cmd MigrateCmd) IsValid() bool {
	if cmd.input == nil || len(cmd.input) == 0 || len(cmd.input) <= 3 {
		return false
	}
	return true
}

func (cmd MigrateCmd) Initialize(input []string) EvebotCommand {
	cmd.input = input
	return cmd
}

func (cmd MigrateCmd) Name() string {
	return cmd.name
}

func (cmd MigrateCmd) Help() *bothelp.Help {
	return bothelp.New(
		bothelp.HeaderOpt(cmd.summary.String()),
		bothelp.UsageOpt(cmd.usage.String()),
		bothelp.ArgsOpt(cmd.optionalArgs.String()),
		bothelp.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd MigrateCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd MigrateCmd) AdditionalArgs() (botargs.Args, error) {
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
