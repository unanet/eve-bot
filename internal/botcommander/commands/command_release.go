package commands

import (
	"fmt"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/help"
	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
)

func NewReleaseCommand(cmdFields []string, channel, user string) EvebotCommand {
	return defaultReleaseCommand(cmdFields, channel, user)
}

type ReleaseCmd struct {
	baseCommand
}

// @evebot release unanet-analytics:20.2 from int to prod
func defaultReleaseCommand(cmdFields []string, channel, user string) ReleaseCmd {
	cmd := ReleaseCmd{baseCommand{
		input:       cmdFields,
		chatDetails: ChatDetails{User: user, Channel: channel},
		name:        "release",
		summary:     "The `release` command is used to release artifacts from/to feeds",
		usage: help.Usage{
			"release {{ artifact }}:{{ optional_version }} from {{ required_feed }}",
			"release {{ artifact }}:{{ optional_version }} from {{ required_feed }} to {{ optional_feed }}",
		},
		examples: help.Examples{
			"release unanet-analytics from int",
			"release unanet-app:20.3 from int",
			"release unanet-analytics:20.2.5 from int to prod",
			"release unanet-analytics:20.2.5.43 from prod to int",
		},
		apiOptions:  make(CommandOptions),
		inputBounds: InputLengthBounds{Min: 4, Max: 6},
	}}
	cmd.resolveValues()
	return cmd
}

func (cmd ReleaseCmd) IsAuthorized(allowedChannelMap map[string]interface{}, fn chatChannelInfo) bool {
	return validChannelAuthCheck(cmd.chatDetails.Channel, allowedChannelMap, fn)
}

func (cmd ReleaseCmd) APIOptions() CommandOptions {
	return cmd.apiOptions
}

func (cmd ReleaseCmd) ChatInfo() ChatDetails {
	return cmd.chatDetails
}

func (cmd ReleaseCmd) AckMsg() (string, bool) {
	return baseAckMsg(cmd, cmd.input)
}

func (cmd ReleaseCmd) IsValid() bool {
	return cmd.ValidInputLength()
}

func (cmd ReleaseCmd) ErrMsg() string {
	return baseErrMsg(cmd.errs)
}

func (cmd ReleaseCmd) Name() string {
	return cmd.name
}

func (cmd ReleaseCmd) Help() *help.Help {
	return help.New(
		help.HeaderOpt(cmd.summary.String()),
		help.UsageOpt(cmd.usage.String()),
		help.ExamplesOpt(cmd.examples.String()),
	)
}

func (cmd ReleaseCmd) IsHelpRequest() bool {
	return isHelpRequest(cmd.input, cmd.name)
}

func (cmd *ReleaseCmd) resolveValues() {
	if cmd.ValidInputLength() == false {
		cmd.errs = append(cmd.errs, fmt.Errorf("invalid release command: %v", cmd.input))
		return
	}

	if strings.Contains(cmd.input[1], ":") {
		artifactKV := strings.Split(cmd.input[1], ":")
		cmd.apiOptions[params.ArtifactName] = artifactKV[0]
		cmd.apiOptions[params.ArtifactVersionName] = artifactKV[1]
	} else {
		cmd.apiOptions[params.ArtifactName] = cmd.input[1]
		cmd.apiOptions[params.ArtifactVersionName] = ""
	}

	cmd.apiOptions[params.FromFeedName] = cmd.input[3]

	switch len(cmd.input) {
	case cmd.inputBounds.Min:
		cmd.apiOptions[params.ToFeedName] = ""
		return
	case cmd.inputBounds.Max:
		cmd.apiOptions[params.ToFeedName] = cmd.input[5]
		return
	}

	return
}
