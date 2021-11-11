package commands

import (
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
)

type releaseArtifactCmd struct {
	baseCommand
}

const (
	// ReleaseArtifactCmdName is the ID/Key for the ReleaseArtifactCmd
	ReleaseArtifactCmdName = "release-artifact"
)

var (
	releaseArtifactCmdHelpSummary = help.Summary("The `release-artifact` command is used to release artifacts from/to feeds")
	releaseArtifactCmdHelpUsage   = help.Usage{
		"release-artifact {{ artifact }}:{{ optional_version }} from {{ required_feed }}",
		"release-artifact {{ artifact }}:{{ optional_version }} from {{ required_feed }} to {{ optional_feed }}",
	}
	releaseArtifactCmdHelpExample = help.Examples{
		"release-artifact api from int",
		"release-artifact api:1.3 from int",
		"release-artifact billing:1.2.4 from int to prod",
		"release-artifact billing:1.2.4 from prod to int",
	}
)

// NewReleaseArtifactCommand creates a New ReleaseArtifactCmd that implements the EvebotCommand interface
func NewReleaseArtifactCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := releaseArtifactCmd{baseCommand{
		input: cmdFields,
		info: ChatInfo{
			User:          user,
			Channel:       channel,
			CommandName:   ReleaseArtifactCmdName,
			IsHelpRequest: isHelpCmd(cmdFields, ReleaseArtifactCmdName),
		},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 4, Max: 6},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
func (cmd releaseArtifactCmd) AckMsg() (string, bool) {
	return cmd.BaseAckMsg(help.New(
		help.HeaderOpt(releaseArtifactCmdHelpSummary.String()),
		help.UsageOpt(releaseArtifactCmdHelpUsage.String()),
		help.ExamplesOpt(releaseArtifactCmdHelpExample.String()),
	).String())
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd releaseArtifactCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd releaseArtifactCmd) Info() ChatInfo {
	return cmd.info
}

func (cmd *releaseArtifactCmd) resolveDynamicOptions() {
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
