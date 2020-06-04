package botcommands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"

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
			"migrate current in una-int databases=unanetbi dryrun=true",
			"migrate current in una-int databases=unanetbi dryrun=true force=true",
			"migrate current in una-int databases=unanetbi,unaneta dryrun=true force=true",
		},
		async:               true,
		optionalArgs:        botargs.Args{botargs.DefaultDryrunArg(), botargs.DefaultForceArg(), botargs.DefaultDatabasesArg()},
		requiredParams:      botparams.Params{botparams.DefaultNamespace(), botparams.DefaultEnvironment()},
		apiOptions:          make(map[string]interface{}),
		requiredInputLength: 4,
	}}
}

func (cmd MigrateCmd) EveReqObj(cbURL, user string) interface{} {
	return eveapi.DeploymentPlanOptions{
		Artifacts:        extractArtifactsOpt(cmd.apiOptions),
		ForceDeploy:      extractForceDeployOpt(cmd.apiOptions),
		User:             user,
		DryRun:           extractDryrunOpt(cmd.apiOptions),
		CallbackURL:      cbURL,
		Environment:      extractEnvironmentOpt(cmd.apiOptions),
		NamespaceAliases: extractNSOpt(cmd.apiOptions),
		Messages:         nil,
		Type:             "migration",
	}
}

func (cmd MigrateCmd) AckMsg(userID string) string {
	return baseAckMsg(cmd, userID, cmd.input)
}

func (cmd MigrateCmd) IsValid() bool {
	if baseIsValid(cmd.input) && len(cmd.input) >= cmd.requiredInputLength {
		return true
	}
	return false
}

func (cmd MigrateCmd) MakeAsyncReq() bool {
	if cmd.IsHelpRequest() || cmd.IsValid() == false || len(cmd.errs) > 0 {
		return false
	}
	return cmd.async
}

func (cmd MigrateCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
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

func (cmd *MigrateCmd) resolveParams() {

	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command params: %v", cmd.input))
		return
	}

	cmd.apiOptions[botparams.NamespaceName] = cmd.input[1]
	cmd.apiOptions[botparams.EnvironmentName] = cmd.input[3]

	return
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
			if suppliedArg := botargs.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.apiOptions[suppliedArg.Name()] = suppliedArg.Value()
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", cmd.input))
			}
		}
	}

	return
}
