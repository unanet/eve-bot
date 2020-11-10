package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/args"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type migrateCmd struct {
	baseCommand
}

const (
	// MigrateCmdName is the ID/Key for the Migrate Command
	MigrateCmdName = "migrate"
)

var (
	migrateCmdHelpSummary = help.Summary("The `migrate` command is used to migrate databases by *namespace* and *environment*")
	migrateCmdHelpUsage   = help.Usage{
		"migrate {{ namespace }} in {{ environment }}",
		"migrate {{ namespace }} in {{ environment }} databases={{ database_type }}",
		"migrate {{ namespace }} in {{ environment }} databases={{ database_type:version }}",
		"migrate {{ namespace }} in {{ environment }} databases={{ database_type:version,database_type:version }} dryrun={{ true }}",
		"migrate {{ namespace }} in {{ environment }} databases={{ database_type:version,database_type }} dryrun={{ true }} force={{ true }}",
	}
	migrateCmdHelpExample = help.Examples{
		"migrate current in qa",
		"migrate current in una-int databases=unanetbi dryrun=true",
		"migrate current in una-int databases=unanetd:20.2 dryrun=true force=true",
		"migrate current in una-int databases=unanetbi,unaneta dryrun=true force=true",
	}
)

// NewMigrateCommand creates a New MigrateCmd that implements the EvebotCommand interface
func NewMigrateCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := migrateCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: MigrateCmdName},
		arguments:  args.Args{args.DefaultDryrunArg(), args.DefaultForceArg(), args.DefaultDatabasesArg()},
		parameters: params.Params{params.DefaultNamespace(), params.DefaultEnvironment()},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 4, Max: 7},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd migrateCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(migrateCmdHelpSummary.String()),
		help.UsageOpt(migrateCmdHelpUsage.String()),
		help.ArgsOpt(cmd.arguments.String()),
		help.ExamplesOpt(migrateCmdHelpExample.String()),
	).String())
}

// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd migrateCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return cmd.IsHelpRequest() || validChannelAuthCheck(cmd.info.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.opts)
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd migrateCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd migrateCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd migrateCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid command params: %v", cmd.input))
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
				cmd.errs = append(cmd.errs, fmt.Errorf("invalid additional arg: %v", cmd.input))
			}
		}
	}
}
