package commands

import (
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
)

type releaseCmd struct {
	baseCommand
}

const (
	// ReleaseCmdName is the ID/Key for the ReleaseCmd
	ReleaseCmdName = "release"
)

var (
	releaseCmdHelpSummary = help.Summary("The `release` command is used to release artifacts from/to feeds")
	releaseCmdHelpUsage   = help.Usage{
		"release {{ artifact }}:{{ optional_version }} from {{ required_feed }}",
		"release {{ artifact }}:{{ optional_version }} from {{ required_feed }} to {{ optional_feed }}",
	}
	releaseCmdHelpExample = help.Examples{
		"release api from int",
		"release api:1.3 from int",
		"release billing:1.2.4 from int to prod",
		"release billing:1.2.4 from prod to int",
	}
)

// NewReleaseCommand creates a New ReleaseCmd that implements the EvebotCommand interface
func NewReleaseCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := releaseCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: ReleaseCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 4, Max: 6},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd releaseCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(releaseCmdHelpSummary.String()),
		help.UsageOpt(releaseCmdHelpUsage.String()),
		help.ExamplesOpt(releaseCmdHelpExample.String()),
	).String())
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd releaseCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd releaseCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *releaseCmd) resolveDynamicOptions() {
	cmd.verifyInput()
	if len(cmd.errs) > 0 {
		return
	}

	if strings.Contains(cmd.input[1], ":") {
		artifactKV := strings.Split(cmd.input[1], ":")
		cmd.opts[params.ArtifactName] = artifactKV[0]
		cmd.opts[params.ArtifactVersionName] = artifactKV[1]
	} else {
		cmd.opts[params.ArtifactName] = cmd.input[1]
		cmd.opts[params.ArtifactVersionName] = ""
	}

	cmd.opts[params.FromFeedName] = cmd.input[3]

	switch len(cmd.input) {
	case cmd.bounds.Min:
		cmd.opts[params.ToFeedName] = ""
		return
	case cmd.bounds.Max:
		cmd.opts[params.ToFeedName] = cmd.input[5]
		return
	}
}
