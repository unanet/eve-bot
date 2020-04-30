package botcommands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
)

func NewMigrateCommand(cmdFields []string) EvebotCommand {
	cmd := defaultMigrateCommand()
	cmd.input = cmdFields
	cmd.resolveParams()
	cmd.resolveArgs()
	return cmd
}

type MigrateCmd struct {
	baseCommand
}

// @evebot migrate current in qa
func defaultMigrateCommand() MigrateCmd {
	return MigrateCmd{baseCommand{
		name:    "migrate",
		summary: "The `migrate` command is used to migrate databases by *namespace* and *environment*",
		usage: bothelp.Usage{
			"migrate {{ namespace }} in {{ environment }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ database_type }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ database_type,database_type }} dryrun={{ true }}",
			"migrate {{ namespace }} in {{ environment }} databases={{ database_type,database_type }} dryrun={{ true }} force={{ true }}",
		},
		examples: bothelp.Examples{
			"migrate current in qa",
			"migrate current in qa databases=infocus dryrun=true",
			"migrate current in qa databases=infocus dryrun=true force=true",
			"migrate current in qa databases=infocus,cloud-support dryrun=true force=true",
		},
		async:          true,
		optionalArgs:   botargs.Args{botargs.DefaultDryrunArg(), botargs.DefaultForceArg(), botargs.DefaultDatabasesArg()},
		suppliedArgs:   botargs.Args{},
		requiredParams: botparams.Params{botparams.DefaultNamespace(), botparams.DefaultEnvironment()},
		suppliedParams: botparams.Params{},
	}}
}

func (cmd *MigrateCmd) resolveParams() {
	if len(cmd.suppliedParams) > 0 {
		return
	}
	if len(cmd.input) < 4 {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command params: %v", cmd.input))
		return
	}
	cmd.suppliedParams = append(
		cmd.suppliedParams,
		botparams.NewNamespaceParam(cmd.input[1]),
		botparams.NewEnvironmentParam(cmd.input[3]),
	)
	return
}

func (cmd *MigrateCmd) resolveArgs() {
	// if we've already calculated the args, use them
	if len(cmd.suppliedArgs) > 0 {
		return
	}

	// haven't calculated the args and no need since they weren't supplied
	if len(cmd.input) < 4 {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command args: %v", cmd.input))
		return
	}

	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if suppliedArg := botargs.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.suppliedArgs = append(cmd.suppliedArgs, suppliedArg)
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", cmd.input))
			}
		}
	}

	return
}

func (cmd MigrateCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd MigrateCmd) AckMsg(userID string) string {
	return baseAckMsg(cmd, userID, cmd.input)
}

func (cmd MigrateCmd) MakeAsyncReq() bool {
	if cmd.IsHelpRequest() || cmd.IsValid() == false || len(cmd.errs) > 0 {
		return false
	}
	return cmd.async
}

func (cmd MigrateCmd) IsValid() bool {
	if baseIsValid(cmd.input) && len(cmd.input) >= 4 {
		return true
	}
	return false
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
