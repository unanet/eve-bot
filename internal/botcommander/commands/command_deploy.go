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

func NewDeployCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := DeployCmd{baseCommand{
		input:       cmdFields,
		chatDetails: ChatInfo{User: user, Channel: channel},
		name:        "deploy",
		summary:     "The `deploy` command is used to deploy services to a specific *namespace* and *environment*",
		usage: help.Usage{
			"deploy {{ namespace }} in {{ environment }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }}",
			"deploy {{ namespace }} in {{ environment }} services={{ service_name:service_version,service_name:service_version }} dryrun={{ true }} force={{ true }}",
		},
		examples: help.Examples{
			"deploy current in una-int",
			"deploy current in una-int services=unanetbi dryrun=true",
			"deploy current in una-int services=unanetbi,unaneta dryrun=true force=true",
			"deploy current in una-int services=unanetbi:20.2,unaneta",
		},
		optionalArgs:   args.Args{args.DefaultDryrunArg(), args.DefaultForceArg(), args.DefaultServicesArg()},
		requiredParams: params.Params{params.DefaultNamespace(), params.DefaultEnvironment()},
		apiOptions:     make(CommandOptions),
		inputBounds:    InputLengthBounds{Min: 4, Max: 7},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

func (cmd DeployCmd) Details() CommandDetails {
	return CommandDetails{
		Name:          cmd.name,
		IsValid:       cmd.ValidInputLength(),
		IsHelpRequest: isHelpRequest(cmd.input, cmd.name),
		AckMsgFn:      baseAckMsg(cmd, cmd.input),
		ErrMsgFn:      cmd.BaseErrMsg(),
	}
}

func (cmd DeployCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return validChannelAuthCheck(cmd.chatDetails.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.apiOptions)
}

func (cmd DeployCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd DeployCmd) ChatInfo() ChatInfo {
	return cmd.chatDetails
}

func (cmd DeployCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ArgsOpt(cmd.optionalArgs.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd *DeployCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("resolve cmd params err invalid input: %v", cmd.input))
		return
	}
	cmd.apiOptions[params.NamespaceName] = cmd.input[1]
	cmd.apiOptions[params.EnvironmentName] = cmd.input[3]

	for _, s := range cmd.input[3:] {
		if strings.Contains(s, "=") {
			argKV := strings.Split(s, "=")
			if suppliedArg := args.ResolveArgumentKV(argKV); suppliedArg != nil {
				cmd.apiOptions[suppliedArg.Name()] = suppliedArg.Value()
			} else {
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", argKV))
			}
		}
	}
}
