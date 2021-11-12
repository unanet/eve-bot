package commands

import (
	"fmt"

	"github.com/unanet/eve-bot/internal/botcommander/help"
)

type factory struct {
	Map map[string]func(cmdFields []string, channel, user string) EvebotCommand
}

type Factory interface {
	Items() map[string]func(cmdFields []string, channel, user string) EvebotCommand
	NonHelpExamples() help.Examples
	NonHelpCmds() string
}

func NewFactory() Factory {
	return &factory{
		Map: map[string]func(cmdFields []string, channel string, user string) EvebotCommand{
			helpCmdName:             NewHelpCommand,
			DeployCmdName:           NewDeployCommand,
			ShowCmdName:             NewShowCommand,
			SetCmdName:              NewSetCommand,
			DeleteCmdName:           NewDeleteCommand,
			ReleaseCmdName:          NewReleaseCommand,
			RestartCmdName:          NewRestartCommand,
			RunCmdName:              NewRunCommand,
			AuthCmdName:             NewAuthCommand,
		},
	}
}

func (f *factory) NonHelpCmds() string {
	var result string
	for k, v := range f.Map {
		if k != helpCmdName {
			nonHelpCmd := v([]string{}, "", "")
			result = result + "\n" + nonHelpCmd.Info().CommandName
		}
	}
	return result
}

func (f *factory) NonHelpExamples() help.Examples {
	var results help.Examples
	for k, v := range f.Map {
		if k != helpCmdName {
			nonHelpCmd := v([]string{}, "", "")
			results = append(results, fmt.Sprintf("%s %s", nonHelpCmd.Info().CommandName, helpCmdName))
		}
	}
	return results
}

func (f *factory) Items() map[string]func(cmdFields []string, channel, user string) EvebotCommand {
	return f.Map
}
