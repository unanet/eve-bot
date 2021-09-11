package service

import (
	"strings"

	"github.com/unanet/eve-bot/internal/botcommander/commands"
	"github.com/unanet/eve-bot/internal/botcommander/params"
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
	return result
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
