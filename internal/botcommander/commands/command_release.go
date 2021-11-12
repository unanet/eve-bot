package commands

import (
	"fmt"
	"github.com/unanet/eve-bot/internal/botcommander/help"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"strings"
)

type releaseCmd struct {
	baseCommand
}

const (
	// ReleaseCmdName is the ID/Key for the ReleaseCmd
	ReleaseCmdName = "release"
)

var (
	releaseNamespaceInputLengthBounds = InputLengthBounds{Min: 6, Max: 8}
	releaseCmdHelpSummary = help.Summary("The `release` command is used to release artifacts or namespaces from/to feeds")
	releaseCmdHelpUsage   = help.Usage{
		// Artifact
		"release artifact {{ artifact }}:{{ optional_version }} from {{ required_feed }}",
		"release artifact {{ artifact }}:{{ optional_version }} from {{ required_feed }} to {{ optional_feed }}\n",

		// Namespace
		"release namespace {{ namespace }} {{ environment }} from {{ required_feed }}",
		"release namespace {{ namespace }} {{ environment }} from {{ required_feed }} to {{ optional_feed }}",
	}
	releaseCmdHelpExample = help.Examples{
		// Artifact
		"release artifact api from int",
		"release artifact api:1.3 from int",
		"release artifact billing:1.2.4 from int to prod",
		"release artifact billing:1.2.4 from prod to int\n",

		// Namespace
		"release namespace current una-int from int",
		"release namespace current una-int from int to stage",
	}
)

// NewReleaseCommand creates a New ReleaseCmd that implements the EvebotCommand interface
func NewReleaseCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := releaseCmd{baseCommand{
		input: cmdFields,
		info: ChatInfo{
			User:          user,
			Channel:       channel,
			CommandName:   ReleaseCmdName,
			IsHelpRequest: isHelpCmd(cmdFields, ReleaseCmdName),
		},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 5, Max: 7},
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

	// Since we combine two ways to handle the same command, we need a little bit of a "hack" to handle
	// the two different lengths
	if len(cmd.input) > 1 && cmd.input[1] == "namespace" {
		cmd.bounds = releaseNamespaceInputLengthBounds
	}

	cmd.verifyInput()
	if len(cmd.errs) > 0 {
		return
	}

	var toBounds int

	switch cmd.input[1] {
		case "artifact":
			cmd.opts["commandType"] = "artifact"

			cmd.resolveDynamicOptionsForArtifact()
			toBounds = 6
		case "namespace":
			cmd.opts["commandType"] = "namespace"

			cmd.resolveDynamicOptionsForNamespace()
			toBounds = 7
		default:
			cmd.errs = append(cmd.errs, fmt.Errorf(
					"unable to determine release type ( artifact or namespace ). " +
							"This has recently changed, please double check with the following examples:\n\n%s",
							releaseCmdHelpUsage,
						),
					)
	}

	switch len(cmd.input) {
	case cmd.bounds.Min:
		cmd.opts[params.ToFeedName] = ""
		return
	case cmd.bounds.Max:
		cmd.opts[params.ToFeedName] = cmd.input[toBounds]
		return
	}
}

func (cmd *releaseCmd) resolveDynamicOptionsForArtifact() {
	if strings.Contains(cmd.input[2], ":") {
		artifactKV := strings.Split(cmd.input[2], ":")
		cmd.opts[params.ArtifactName] = artifactKV[0]
		cmd.opts[params.ArtifactVersionName] = artifactKV[1]
	} else {
		cmd.opts[params.ArtifactName] = cmd.input[2]
		cmd.opts[params.ArtifactVersionName] = ""
	}

	cmd.opts[params.FromFeedName] = cmd.input[4]
}

func (cmd *releaseCmd) resolveDynamicOptionsForNamespace() {
	cmd.opts[params.NamespaceName] = cmd.input[2]
	cmd.opts[params.EnvironmentName] = cmd.input[3]

	cmd.opts[params.FromFeedName] = cmd.input[5]
}