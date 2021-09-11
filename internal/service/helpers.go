package service

import (
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
	"github.com/unanet/go/pkg/log"
	"go.uber.org/zap"
)

func extractIsAdminRole(input interface{}) bool {
	if roles, ok := input.([]interface{}); ok {
		for _, role := range roles {
			if strings.Contains(strings.ToLower(role.(string)), "admin") {
				return true
			}
		}

	}
	return false
}

func extractClaimMap(input interface{}) map[string]bool {
	result := make(map[string]bool)
	if v, ok := input.([]interface{}); ok {
		for _, param := range v {
			result[param.(string)] = true
		}

	}
	log.Logger.Warn("invalid type on incoming claim slice", zap.Any("input", input), zap.Reflect("type", input))
	return result
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
