package commands

import (
	"context"
	"strings"

	"gitlab.unanet.io/devops/eve-bot/internal/botcommander/params"
	"gitlab.unanet.io/devops/eve/pkg/log"
	"go.uber.org/zap"
)

// validChannelAuthCheck validates/confirm if the incoming channel matches one of the "approved" channels
// approved channels configured via Environment Variable: EVEBOT_SLACK_CHANNELS_AUTH
func validChannelAuthCheck(channel string, channelMap map[string]interface{}, fn chatChannelInfoFn) bool {
	incomingChannelInfo, err := fn(context.TODO(), channel)
	if err != nil {
		log.Logger.Error("failed to get channel info", zap.Error(err))
		return false
	}

	// Coming from an Elevated/Approved Channel
	// let them pass
	if _, ok := channelMap[incomingChannelInfo.Name]; ok {
		return true
	}
	return false
}

// lowerEnvAuthCheck checks if the incoming command is against a lowe-environment (int/qa/dev)
func lowerEnvAuthCheck(options CommandOptions) bool {
	if options == nil {
		return false
	}

	if env, ok := options[params.EnvironmentName].(string); ok {
		// Let's see if they are performing an action to something in the lower environments (int,qa,dev)
		// Most actions can be taken against resources in the lower environments
		// the only action that can't is the `release` command
		switch {
		case strings.Contains(env, "int"), strings.Contains(env, "qa"), strings.Contains(env, "dev"):
			return true
		}
		return false
	}

	log.Logger.Warn("environment not set")
	return false
}
