package service

import (
	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
	"strings"
)

func extractChannelMap(input string) map[string]interface{} {
	chanMap := make(map[string]interface{})
	for _, c := range strings.Split(input, ",") {
		chanMap[c] = true
	}
	return chanMap
}

func extractClaimSlice(input interface{}) []string {
	if v, ok := input.([]interface{}); ok {
		var paramSlice []string
		for _, param := range v {
			paramSlice = append(paramSlice, param.(string))
		}
		return paramSlice
	}
	log.Logger.Warn("invalid type on incoming claim slice", zap.Any("input", input), zap.Reflect("type", input))
	return []string{}
}

func extractEnv(options commands.CommandOptions) string {
	if options == nil {
		return ""
	}
	if env, ok := options[params.EnvironmentName].(string); ok {
		return env
	}
	return ""
}
