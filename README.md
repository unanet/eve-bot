# eve-bot

## Summary

This is the `eve-bot` ChatOps service. It is responsible for handling all communication between the User (Slack) and the Backend Pipeline API.

## Building

This project utilizes `Make` for the build process. The `Makefile` calls shell scripts for various tasks, and acts like a think wrapper

### Build Standard Golang Binary (Go is required)

1. `make build`

### Build Docker Image (Docker is required)

1. `make docker`

## Running

This application uses sane defaults for most of the config, but there are some required secrets that need to be set as  `Environment Variables`. **All application config use EnvVars.**

### Application Environment Variables

```bash
EVEBOT_API_SERVICENAME
EVEBOT_API_PORT
EVEBOT_API_METRICSPORT
EVEBOT_API_ALLOWEDMETHODS
EVEBOT_API_ALLOWEDORIGINS
EVEBOT_API_ALLOWEDHEADERS
EVEBOT_LOGGER_LEVEL
EVEBOT_LOGGER_ENCODING
EVEBOT_LOGGER_OUTPUTPATHS
EVEBOT_LOGGER_ERROUTPUTPATHS
EVEBOT_SLACK_SIGNING_SECRET
EVEBOT_SLACK_VERIFICATION_TOKEN
EVEBOT_SLACK_BOT_OAUTH
EVEBOT_SLACK_OAUTH
```

All environment variables are represented as strings. The list/slice variables need to be set as a comma separated string:

```bash
export EVEBOT_API_ALLOWEDMETHODS="GET,HEAD,POST,PUT,OPTIONS,DELETE"
export EVEBOT_API_ALLOWEDORIGINS="localhost,evebot.unanet.io,api.slack.com"
export EVEBOT_API_ALLOWEDHEADERS="X-Requested-With,X-Request-ID,jaeger-debug-id,Content-Type,X-Slack-Signature,X-Slack-Request-Timestamp"
export EVEBOT_LOGGER_OUTPUTPATHS="stdout,/tmp/evebot.logs"
export EVEBOT_LOGGER_ERROUTPUTPATHS="stderr"
```

### Required Environment Variables

```bash
export EVEBOT_SLACK_SIGNING_SECRET=`vault kv get --format=json kv/devops/evebot | jq .data.data.EVEBOT_SLACK_SIGNING_SECRET`
export EVEBOT_SLACK_VERIFICATION_TOKEN=`vault kv get --format=json kv/devops/evebot | jq .data.data.EVEBOT_SLACK_VERIFICATION_TOKEN`
export EVEBOT_SLACK_BOT_OAUTH=`vault kv get --format=json kv/devops/evebot | jq .data.data.EVEBOT_SLACK_BOT_OAUTH`
export EVEBOT_SLACK_OAUTH=`vault kv get --format=json kv/devops/evebot | jq .data.data.EVEBOT_SLACK_OAUTH`
```

These secrets are required for the application to run. The source of truth is Slack, but we store them in Vault for safe keeping. If we need to roll the secrets (generate knew ones) that should be done through the Slack UI, and then pushed up to Vault.

### Slack Links

[Slack OAuth Tokens](https://api.slack.com/apps/A011B3L27P1/oauth)

[Slack Event Subscriptions](https://api.slack.com/apps/A011B3L27P1/event-subscriptions)

[Slack App Creds](https://api.slack.com/apps/A011B3L27P1/general?)