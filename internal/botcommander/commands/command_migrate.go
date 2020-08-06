package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

func NewMigrateCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := MigrateCmd{baseCommand{
		input:       cmdFields,
		chatDetails: ChatInfo{User: user, Channel: channel},
		name:        "migrate",
		summary:     "The `migrate` command is used to migrate databases by *namespace* and *environment*",
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
		optionalArgs:   args.Args{args.DefaultDryrunArg(), args.DefaultForceArg(), args.DefaultDatabasesArg()},
		requiredParams: params.Params{params.DefaultNamespace(), params.DefaultEnvironment()},
		apiOptions:     make(CommandOptions),
		inputBounds:    InputLengthBounds{Min: 4, Max: 7},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

type MigrateCmd struct {
	baseCommand
}

func (cmd MigrateCmd) Details() CommandDetails {
	return CommandDetails{
		Name:          cmd.name,
		IsValid:       cmd.ValidInputLength(),
		IsHelpRequest: isHelpRequest(cmd.input, cmd.name),
		AckMsgFn:      baseAckMsg(cmd, cmd.input),
		ErrMsgFn:      cmd.BaseErrMsg(),
	}
}

func (cmd MigrateCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return validChannelAuthCheck(cmd.chatDetails.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.apiOptions)
}

func (cmd MigrateCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd MigrateCmd) ChatInfo() ChatInfo {
	return cmd.chatDetails
}

func (cmd MigrateCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ArgsOpt(cmd.optionalArgs.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd MigrateCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command params: %v", cmd.input))
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
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", cmd.input))
			}
		}
	}
}
