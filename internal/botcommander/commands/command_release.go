package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

type ReleaseCmd struct {
	baseCommand
}

const (
	ReleaseCmdName = "release"
)

var (
	releaseCmdHelpSummary = help.Summary("The `release` command is used to release artifacts from/to feeds")
	releaseCmdHelpUsage   = help.Usage{
		"release {{ artifact }}:{{ optional_version }} from {{ required_feed }}",
		"release {{ artifact }}:{{ optional_version }} from {{ required_feed }} to {{ optional_feed }}",
	}
	releaseCmdHelpExample = help.Examples{
		"release unanet-analytics from int",
		"release unanet-app:20.3 from int",
		"release unanet-analytics:20.2.5 from int to prod",
		"release unanet-analytics:20.2.5.43 from prod to int",
	}
)

func NewReleaseCommand(cmdFields []string, channel, user string) EvebotCommand {
	cmd := ReleaseCmd{baseCommand{
		input:  cmdFields,
		info:   ChatInfo{User: user, Channel: channel, CommandName: ReleaseCmdName},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: 4, Max: 6},
	}}
	cmd.resolveDynamicOptions()
	return cmd
}

func (cmd ReleaseCmd) AckMsg() (string, bool) {

	helpMsg := help.New(
		help.HeaderOpt(releaseCmdHelpSummary.String()),
		help.UsageOpt(releaseCmdHelpUsage.String()),
		help.ExamplesOpt(releaseCmdHelpExample.String()),
	).String()

	return cmd.BaseAckMsg(helpMsg)
}

func (cmd ReleaseCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	return validChannelAuthCheck(cmd.info.Channel, allowedChannelMap, fn)
}

func (cmd ReleaseCmd) DynamicOptions() CommandOptions {
	return cmd.opts
}

func (cmd ReleaseCmd) ChatInfo() ChatInfo {
	return cmd.info
}

func (cmd *ReleaseCmd) resolveDynamicOptions() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid release command: %v", cmd.input))
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

	return
}
