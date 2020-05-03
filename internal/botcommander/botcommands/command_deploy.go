package botcommands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botargs"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/bothelp"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/botparams"
	"gitlab.unanet.io/devops/eve-bot/internal/eveapi"
)

func NewDeployCommand(cmdFields []string) EvebotCommand {
	cmd := defaultDeployCommand()
	cmd.input = cmdFields
	cmd.resolveParams()
	cmd.resolveArgs()
	return cmd
}

type DeployCmd struct {
	baseCommand
}

func defaultDeployCommand() DeployCmd {
	return DeployCmd{baseCommand{
		name:    "deploy",
		summary: "The `deploy` command is used to deploy services to a specific *namespace* and *environment*",
		usage: bothelp.Usage{
			"deploy {{ namespace }} in {{ environment }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		examples: bothelp.Examples{
			"deploy current in qa",
			"deploy current in qa services=infocus-cloud-client:2020.1 dryrun=true",
			"deploy current in qa services=infocus-cloud-client:2020.2.232,infocus-proxy:2020.2.199 dryrun=true force=true",
			"deploy current in qa services=infocus-cloud-client,infocus-proxy",
		},
		async:          true,
		optionalArgs:   botargs.Args{botargs.DefaultDryrunArg(), botargs.DefaultForceArg(), botargs.DefaultServicesArg()},
		requiredParams: botparams.Params{botparams.DefaultNamespace(), botparams.DefaultEnvironment()},
		apiOptions:     make(map[string]interface{}),
	}}
}

//Artifacts   ArtifactDefinitions `json:"artifacts"`
//ForceDeploy bool                `json:"force_deploy"`
//DryRun      bool                `json:"dry_run"`
//CallbackURL string              `json:"callback_url"`
//Environment string              `json:"environment"`
//Namespaces  []string            `json:"namespaces,omitempty"`
//Messages    []string            `json:"messages,omitempty"`
//Type        string              `json:"type"`
func (cmd DeployCmd) EveReqObj() interface{} {

	opts := eveapi.DeploymentPlanOptions{
		CallbackURL: "config.Values().CallbackURL",
		Type:        "application",
	}

	if val, ok := cmd.apiOptions[botargs.ServicesName]; ok {
		opts.Artifacts = val.(eveapi.ArtifactDefinitions)
	}

	if val, ok := cmd.apiOptions[botargs.ForceDeployName]; ok {
		opts.ForceDeploy = val.(bool)
	}

	if val, ok := cmd.apiOptions[botargs.DryrunName]; ok {
		opts.DryRun = val.(bool)
	}

	if val, ok := cmd.apiOptions[botparams.EnvironmentName]; ok {
		opts.Environment = val.(string)
	}

	if val, ok := cmd.apiOptions[botparams.NamespaceName]; ok {
		opts.Namespaces = []string{val.(string)}
	}

	return opts
}

func (cmd DeployCmd) AckMsg(userID string) string {
	return baseAckMsg(cmd, userID, cmd.input)
}

func (cmd DeployCmd) IsValid() bool {
	if baseIsValid(cmd.input) && len(cmd.input) >= 4 {
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
	if len(cmd.input) < 4 {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command params: %v", cmd.input))
		return
	}
	cmd.apiOptions["namespace"] = cmd.input[1]
	cmd.apiOptions["environment"] = cmd.input[3]

	return
}

// resolveArgs attempts to resolve the input argument
// be sure and use a pointer receiver here since we are modifying the receiver object
func (cmd *DeployCmd) resolveArgs() {
	// haven't calculated the args and no need since they weren't supplied
	if len(cmd.input) < 4 {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command params: %v", cmd.input))
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

	return
}
