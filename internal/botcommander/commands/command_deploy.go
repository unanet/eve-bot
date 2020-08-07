package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type DeployCmd struct {
	baseCommand
}

const (
	DeployCmdName = "deploy"
)

var (
	deployCmdHelpSummary = help.Summary("The `deploy` command is used to deploy services to a specific *namespace* and *environment*")
	deployCmdHelpUsage   = help.Usage{
		"deploy {{ namespace }} in {{ environment }}",
		"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
		"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
		"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
	}
	deployCmdHelpExample = help.Examples{
		"deploy current in una-int",
		"deploy current in una-int services=unanetbi dryrun=true",
		"deploy current in una-int services=unanetbi,unaneta dryrun=true force=true",
		"deploy current in una-int services=unanetbi:20.2,unaneta",
	}
)

func NewDeployCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := DeployCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: DeployCmdName},
		arguments:  args.Args{args.DefaultDryrunArg(), args.DefaultForceArg(), args.DefaultServicesArg()},
		parameters: params.Params{params.DefaultNamespace(), params.DefaultEnvironment()},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 4, Max: 7},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

func (cmd DeployCmd) AckMsg() (string, bool) {

	helpMsg := help.New(
		help.HeaderOpt(deployCmdHelpSummary.String()),
		help.UsageOpt(deployCmdHelpUsage.String()),
		help.ArgsOpt(cmd.arguments.String()),
		help.ExamplesOpt(deployCmdHelpExample.String()),
	).String()

	return cmd.BaseAckMsg(helpMsg)
}

func (cmd DeployCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return validChannelAuthCheck(cmd.info.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.opts)
}

func (cmd DeployCmd) DynamicOptions() CommandOptions {
	return cmd.opts
}

func (cmd DeployCmd) ChatInfo() ChatInfo {
	return cmd.info
}

func (cmd *DeployCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("resolve cmd params err invalid input: %v", cmd.input))
		return
	}
	cmd.opts[params.NamespaceName] = cmd.input[1]
	cmd.opts[params.EnvironmentName] = cmd.input[3]

	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if suppliedArg := args.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.opts[suppliedArg.Name()] = suppliedArg.Value()
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", argKV))
			}
		}
	}
}
