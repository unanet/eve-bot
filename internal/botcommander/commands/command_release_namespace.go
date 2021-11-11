package commands

import (
	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
)

type releaseNamespaceCmd struct {
	baseCommand
}

const (
	// ReleaseNamespaceCmdName is the ID/Key for the ReleaseNamespaceCmd
	ReleaseNamespaceCmdName = "release-namespace"
)

var (
	releaseNamespaceCmdHelpSummary = help.Summary("The `release-namespace` command is used to release an entire namespace's artifacts from/to feeds")
	releaseNamespaceCmdHelpUsage   = help.Usage{
		"release-namespace {{ namespace }} {{ environment }} from {{ required_feed }}",
		"release-namespace {{ namespace }} {{ environment }} from {{ required_feed }} to {{ optional_feed }}",
	}
	releaseNamespaceHelpExample = help.Examples{
		"release-namespace current una-int from int",
		"release-namespace current una-int from int to prod",
		"release-namespace current una-int from prod to int",
	}
)

// NewReleaseNamespaceCommand creates a New ReleaseArtifactCmd that implements the EvebotCommand interface
func NewReleaseNamespaceCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := releaseNamespaceCmd{baseCommand{
		input: cmdFields,
		info: ChatInfo{
			User:          user,
			Channel:       channel,
			CommandName:   ReleaseNamespaceCmdName,
			IsHelpRequest: isHelpCmd(cmdFields, ReleaseNamespaceCmdName),
		},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 5, Max: 7},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd releaseNamespaceCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(releaseNamespaceCmdHelpSummary.String()),
		help.UsageOpt(releaseNamespaceCmdHelpUsage.String()),
		help.ExamplesOpt(releaseNamespaceHelpExample.String()),
	).String())
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd releaseNamespaceCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd releaseNamespaceCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *releaseNamespaceCmd) resolveDynamicOptions() {
	cmd.verifyInput()
	if len(cmd.errs) > 0 {
		return
	}

	cmd.opts[params.NamespaceName] = cmd.input[1]
	cmd.opts[params.EnvironmentName] = cmd.input[2]
	cmd.opts[params.FromFeedName] = cmd.input[4]

	switch len(cmd.input) {
	case cmd.bounds.Min:
		cmd.opts[params.ToFeedName] = ""
		return
	case cmd.bounds.Max:
		cmd.opts[params.ToFeedName] = cmd.input[6]
		return
	}
}
