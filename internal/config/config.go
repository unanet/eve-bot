package config

import (
	"log"
	"os"
	"strings"
)

type slackSecrets struct {
	SigningSecret, VerificationToken, BotOAuthToken, OAuthToken string
}

type logger struct {
	Level, Encoding               string
	OutputPaths, ErrorOutputPaths []string
}

type api struct {
	ShutdownTimeoutSecs uint16
	ReadTimeOutSecs     uint16
	WriteTimeOutSecs    uint16
	IdleTimeOutSecs     uint16
	TimeoutSecs         uint16
	ServiceName         string
	Port                string
	MetricsPort         string
	AllowedMethods      []string
	AllowedOrigins      []string
	AllowedHeaders      []string
}

// Config is the main Bot Config
type Config struct {
	API          api
	Logger       logger
	SlackSecrets slackSecrets
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if fallback != "" {
		return fallback
	}

	// The EnvVar isn't set and it's required...time to blow up
	log.Panicf("missing required config: %s", key)
	return ""
}

// Read the config values (iether sane default or required, all EnvVars)
// this should only be called once on app startup
func Read() *Config {
	return &Config{
		API: api{
			ServiceName:         getenv("EVEBOT_API_SERVICENAME", "EveBot"),
			ShutdownTimeoutSecs: 120,
			ReadTimeOutSecs:     5,
			WriteTimeOutSecs:    30,
			IdleTimeOutSecs:     90,
			TimeoutSecs:         30,
			Port:                getenv("EVEBOT_API_PORT", "3000"),
			MetricsPort:         getenv("EVEBOT_API_METRICSPORT", "3001"),
			AllowedMethods:      strings.Split(getenv("EVEBOT_API_ALLOWEDMETHODS", "GET,HEAD,POST,PUT,OPTIONS,DELETE"), ","),
			AllowedOrigins:      strings.Split(getenv("EVEBOT_API_ALLOWEDORIGINS", "*"), ","),
			AllowedHeaders:      strings.Split(getenv("EVEBOT_API_ALLOWEDHEADERS", "*"), ","),
		},
		Logger: logger{
			Level:            getenv("EVEBOT_LOGGER_LEVEL", "debug"),
			Encoding:         getenv("EVEBOT_LOGGER_ENCODING", "json"),
			OutputPaths:      strings.Split(getenv("EVEBOT_LOGGER_OUTPUT_PATHS", "stdout,/tmp/evebot.logs"), ","),
			ErrorOutputPaths: strings.Split(getenv("EVEBOT_LOGGER_ERR_OUTPUT_PATHS", "stderr"), ","),
		},
		SlackSecrets: slackSecrets{
			SigningSecret:     getenv("EVEBOT_SLACK_SIGNING_SECRET", ""),
			VerificationToken: getenv("EVEBOT_SLACK_VERIFICATION_TOKEN", ""),
			BotOAuthToken:     getenv("EVEBOT_SLACK_BOT_OAUTH", ""),
			OAuthToken:        getenv("EVEBOT_SLACK_OAUTH", ""),
		},
	}
}
