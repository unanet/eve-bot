package botcommands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

func NewDeployCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultDeployCommand(cmdFields, channel, user)
}

type DeployCmd struct {
	baseCommand
}

func defaultDeployCommand(cmdFields []string, channel, user string) DeployCmd {
	cmd := DeployCmd{baseCommand{
		input:   cmdFields,
		channel: channel,
		user:    user,
		name:    "deploy",
		summary: "The `deploy` command is used to deploy services to a specific *namespace* and *environment*",
		usage: bothelp.Usage{
			"deploy {{ namespace }} in {{ environment }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		examples: bothelp.Examples{
			"deploy current in una-int",
			"deploy current in una-int services=unanetbi dryrun=true",
			"deploy current in una-int services=unanetbi,unaneta dryrun=true force=true",
			"deploy current in una-int services=unanetbi:20.2,unaneta",
		},
		async:               true,
		optionalArgs:        botargs.Args{botargs.DefaultDryrunArg(), botargs.DefaultForceArg(), botargs.DefaultServicesArg()},
		requiredParams:      botparams.Params{botparams.DefaultNamespace(), botparams.DefaultEnvironment()},
		apiOptions:          make(map[string]interface{}),
		requiredInputLength: 4,
	}}
	cmd.resolveParams()
	cmd.resolveArgs()
	return cmd
}

// EveReqObj hydrates the data needed to make the EveAPI Request for the EveBot Command (deploy)
func (cmd DeployCmd) EveReqObj(user string) interface{} {
	return eveapi.DeploymentPlanOptions{
		Artifacts:        extractArtifactsOpt(cmd.apiOptions),
		ForceDeploy:      extractForceDeployOpt(cmd.apiOptions),
		User:             user,
		DryRun:           extractDryrunOpt(cmd.apiOptions),
		Environment:      extractEnvironmentOpt(cmd.apiOptions),
		NamespaceAliases: extractNSOpt(cmd.apiOptions),
		Messages:         nil,
		Type:             "application",
	}
}

func (cmd DeployCmd) User() string {
	return cmd.user
}

func (cmd DeployCmd) Channel() string {
	return cmd.channel
}

func (cmd DeployCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd DeployCmd) IsValid() bool {
	if baseIsValid(cmd.input) && len(cmd.input) >= cmd.requiredInputLength {
		return true
	}
	return false
}

func (cmd DeployCmd) MakeAsyncReq() bool {
	if cmd.IsHelpRequest() || cmd.IsValid() == false || len(cmd.errs) > 0 {
		return false
	}
	return cmd.async
}

func (cmd DeployCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
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

// resolveParams attempts to resolve the input params
// be sure and use a pointer receiver here since we are modifying the receiver object
func (cmd *DeployCmd) resolveParams() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("resolve cmd params err invalid input: %v", cmd.input))
		return
	}
	cmd.apiOptions[botparams.NamespaceName] = cmd.input[1]
	cmd.apiOptions[botparams.EnvironmentName] = cmd.input[3]
}

// resolveArgs attempts to resolve the input argument
// be sure and use a pointer receiver here since we are modifying the receiver object
func (cmd *DeployCmd) resolveArgs() {
	if len(cmd.input) < cmd.requiredInputLength {
		cmd.errs = append(cmd.errs, fmt.Errorf("resolve cmd args err invalid input: %v", cmd.input))
		return
	}
	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if suppliedArg := botargs.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.apiOptions[suppliedArg.Name()] = suppliedArg.Value()
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", argKV))
			}
		}
	}
}
