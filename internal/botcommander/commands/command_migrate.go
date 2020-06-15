package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

func NewMigrateCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultMigrateCommand(cmdFields, channel, user)
}

type MigrateCmd struct {
	baseCommand
}

// @evebot migrate current in qa
func defaultMigrateCommand(cmdFields []string, channel, user string) MigrateCmd {
	cmd := MigrateCmd{baseCommand{
		input:   cmdFields,
		channel: channel,
		user:    user,
		name:    "migrate",
		summary: "The `migrate` command is used to migrate databases by *namespace* and *environment*",
		usage: help.Usage{
			"migrate {{ namespace }} in {{ environment }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ database_type }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ database_type:version }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ database_type:version,database_type:version }} dryrun={{ true }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ database_type:version,database_type }} dryrun={{ true }} force={{ true }}",
		},
		examples: help.Examples{
			"migrate current in qa",
			"migrate current in una-int databases=unanetbi dryrun=true",
			"migrate current in una-int databases=unanetd:20.2 dryrun=true force=true",
			"migrate current in una-int databases=unanetbi,unaneta dryrun=true force=true",
		},
		optionalArgs:        args.Args{args.DefaultDryrunArg(), args.DefaultForceArg(), args.DefaultDatabasesArg()},
		requiredParams:      params.Params{params.DefaultNamespace(), params.DefaultEnvironment()},
		apiOptions:          make(CommandOptions),
		requiredInputLength: 4,
	}}
	cmd.resolveParams()
	cmd.resolveArgs()
	return cmd
}

func (cmd MigrateCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd MigrateCmd) User() string {
	return cmd.user
}

func (cmd MigrateCmd) Channel() string {
	return cmd.channel
}

func (cmd MigrateCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd MigrateCmd) IsValid() bool {
	if baseIsValid(cmd.input) && len(cmd.input) >= cmd.requiredInputLength {
		return true
	}
	return false
}

func (cmd MigrateCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd MigrateCmd) Name() string {
	return cmd.name
}

func (cmd MigrateCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ArgsOpt(cmd.optionalArgs.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd MigrateCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd *MigrateCmd) resolveParams() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command params: %v", cmd.input))
		return
	}
	cmd.apiOptions[params.NamespaceName] = cmd.input[1]
	cmd.apiOptions[params.EnvironmentName] = cmd.input[3]
}

func (cmd *MigrateCmd) resolveArgs() {
	// haven't calculated the args and no need since they weren't supplied
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command args: %v", cmd.input))
		return
	}
	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if suppliedArg := args.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.apiOptions[suppliedArg.Name()] = suppliedArg.Value()
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", cmd.input))
			}
		}
	}
}
