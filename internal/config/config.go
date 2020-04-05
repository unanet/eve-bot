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

const (
	envVarAPIServiceName         = "EVEBOT_API_SERVICENAME"
	envVarAPIPort                = "EVEBOT_API_PORT"
	envVarAPIMetricPort          = "EVEBOT_API_METRICSPORT"
	envVarAPIAllowedMethods      = "EVEBOT_API_ALLOWEDMETHODS"
	envVarAPIAllowedOrigins      = "EVEBOT_API_ALLOWEDORIGINS"
	envVarAPIAllowedHeaders      = "EVEBOT_API_ALLOWEDHEADERS"
	envVarLoggerLevel            = "EVEBOT_LOGGER_LEVEL"
	envVarLoggerEncoding         = "EVEBOT_LOGGER_ENCODING"
	envVarLoggerOutputPaths      = "EVEBOT_LOGGER_OUTPUTPATHS"
	envVarLoggerErrOutputPaths   = "EVEBOT_LOGGER_ERROUTPUTPATHS"
	envVarSlackSigningSecret     = "EVEBOT_SLACK_SIGNING_SECRET"
	envVarSlackVerificationToken = "EVEBOT_SLACK_VERIFICATION_TOKEN"
	envVarSlackBotOAuth          = "EVEBOT_SLACK_BOT_OAUTH"
	envVarSlackOAuth             = "EVEBOT_SLACK_OAUTH"
)

// Read the config values (iether sane default or required, all EnvVars)
// this should only be called once on app startup
func Read() *Config {
	return &Config{
		API: api{
			ServiceName:         getenv(envVarAPIServiceName, "EveBot"),
			ShutdownTimeoutSecs: 120,
			ReadTimeOutSecs:     5,
			WriteTimeOutSecs:    30,
			IdleTimeOutSecs:     90,
			TimeoutSecs:         30,
			Port:                getenv(envVarAPIPort, "3000"),
			MetricsPort:         getenv(envVarAPIMetricPort, "3001"),
			AllowedMethods:      strings.Split(getenv(envVarAPIAllowedMethods, "GET,HEAD,POST,PUT,OPTIONS,DELETE"), ","),
			AllowedOrigins:      strings.Split(getenv(envVarAPIAllowedOrigins, "*"), ","),
			AllowedHeaders:      strings.Split(getenv(envVarAPIAllowedHeaders, "*"), ","),
		},
		Logger: logger{
			Level:            getenv(envVarLoggerLevel, "debug"),
			Encoding:         getenv(envVarLoggerEncoding, "json"),
			OutputPaths:      strings.Split(getenv(envVarLoggerOutputPaths, "stdout,/tmp/evebot.logs"), ","),
			ErrorOutputPaths: strings.Split(getenv(envVarLoggerErrOutputPaths, "stderr"), ","),
		},
		SlackSecrets: slackSecrets{
			SigningSecret:     getenv(envVarSlackSigningSecret, ""),
			VerificationToken: getenv(envVarSlackVerificationToken, ""),
			BotOAuthToken:     getenv(envVarSlackBotOAuth, ""),
			OAuthToken:        getenv(envVarSlackOAuth, ""),
		},
	}
}
