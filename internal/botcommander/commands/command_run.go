package commands

import (
	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
)

type runCmd struct {
	baseCommand
}

const (
	RunCmdName = "run"
)

var (
	runCmdHelpSummary = help.Summary(
		"The `run` command is used to run a job in a namespace",
	)

	runCmdHelpUsage = help.Usage{
		"run {{ job:version }} in {{ namespace }} {{ environment }} {{key=value}}",
	}

	runCmdHelpExample = help.Examples{
		"run migration in current int key=value key2=value2 keyN=valN",
		"run cool-job:1.2 in current int key=value key2=value2 keyN=valN",
	}
)

// NewRunCommand creates a New RunCmd that implements the EvebotCommand interface
func NewRunCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := runCmd{baseCommand{
		input:      cmdFields,
		info:       ChatInfo{User: user, Channel: channel, CommandName: RunCmdName},
		parameters: params.Params{params.DefaultJob(), params.DefaultNamespace(), params.DefaultEnvironment()},
		opts:       make(CommandOptions),
		bounds:     InputLengthBounds{Min: 5, Max: -1},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd runCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(runCmdHelpSummary.String()),
		help.UsageOpt(runCmdHelpUsage.String()),
		help.ExamplesOpt(runCmdHelpExample.String()),
	).String())
}

// IsAuthorized satisfies the EveBotCommand Interface and checks the auth
func (cmd runCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return cmd.IsHelpRequest() || validChannelAuthCheck(cmd.info.Channel, allowedChannelMap, fn) || lowerEnvAuthCheck(cmd.opts)
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd runCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd runCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *runCmd) resolveDynamicOptions() {
	cmd.verifyInput()
	if len(cmd.errs) > 0 {
		return
	}

	// run {{ job }} in {{ namespace }} {{ environment }}
	// run create-customer in current una-int name=username domain=something.something.somewhere
	cmd.opts[params.JobName] = cmd.input[1]
	cmd.opts[params.NamespaceName] = cmd.input[3]
	cmd.opts[params.EnvironmentName] = cmd.input[4]
	cmd.opts[params.MetadataName] = hydrateMetadataMap(cmd.input[5:])
}
