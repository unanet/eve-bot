package commands

import (
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

func validChannelAuthCheck(channel string, channelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	incomingChannelInfo, err := fn(channel)
	if err != nil {
		log.Logger.Error("failed to get channel info", zap.Error(err))
		return false
	}
	log.Logger.Debug("auth channel check", zap.String("id", incomingChannelInfo.ID), zap.String("name", incomingChannelInfo.Name))

	// Coming from an Elevated/Approved Channel
	// let them pass
	if _, ok := channelMap[incomingChannelInfo.Name]; ok {
		return true
	}
	return false
}

func lowerEnvAuthCheck(options CommandOptions) bool {
	if options == nil {
		return false
	}

	var env string
	var ok bool
	if env, ok = options[params.EnvironmentName].(string); ok == false {
		log.Logger.Warn("environment not set")
		return false
	}

	// Let's see if they are performing an action to something in the lower environments (int,qa,dev)
	// Most actions can be taken against resources in the lower environments
	// the only action that can't is the `release` command
	switch {
	case strings.Contains(env, "int"), strings.Contains(env, "qa"), strings.Contains(env, "dev"):
		return true
	}
	return false
}
