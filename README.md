# eve-bot

[![pipeline status](https://gitlab.unanet.io/devops/eve-bot/badges/master/pipeline.svg)](https://gitlab.unanet.io/devops/eve-bot/-/commits/master)

## Summary

This is the `eve-bot` ChatOps service. It is responsible for handling all communication between the User (Slack) and the Backend Pipeline API.

## Building

run: `make`

## Running

This application uses sane defaults for most of the config, but there are some required secrets that need to be set as `Environment Variables`. **All application config use EnvVars.**

## Environment Variables

```bash
EVE_SERVICE_NAME=eve-bot
EVE_LOG_LEVEL=debug
EVEBOT_PORT=3000
EVEBOT_METRICS_PORT=3001
EVEBOT_SLACK_SIGNING_SECRET=""
EVEBOT_SLACK_VERIFICATION_SECRET=""
EVEBOT_SLACK_USER_OAUTH_ACCESS_TOKEN=""
EVEBOT_SLACK_OAUTH_ACCESS_TOKEN=""
EVEBOT_SLACK_SKIP_VERIFICATION=""
```

These secrets are required for the application to run. The source of truth is Slack, but we store them in Vault for safe keeping.

```bash
export EVEBOT_SLACK_SIGNING_SECRET=`vault kv get --format=json kv/devops/evebot | jq .data.data.EVEBOT_SLACK_SIGNING_SECRET`
export EVEBOT_SLACK_VERIFICATION_TOKEN=`vault kv get --format=json kv/devops/evebot | jq .data.data.EVEBOT_SLACK_VERIFICATION_TOKEN`
export EVEBOT_SLACKBOT_OAUTH_TOKEN=`vault kv get --format=json kv/devops/evebot | jq .data.data.EVEBOT_SLACK_BOT_OAUTH`
export EVEBOT_SLACK_OAUTH_TOKEN=`vault kv get --format=json kv/devops/evebot | jq .data.data.EVEBOT_SLACK_OAUTH`
```

### Slack Links

New secrets should be generated through the Slack UI, and then pushed up to Vault.

[Slack OAuth Tokens](https://api.slack.com/apps/A011B3L27P1/oauth)

[Slack Event Subscriptions](https://api.slack.com/apps/A011B3L27P1/event-subscriptions)

[Slack App Creds](https://api.slack.com/apps/A011B3L27P1/general?)

### Local Dev

To run/develop locally: `docker-compose up` **Note:Still need to setup ngrok to proxy bot/slack request**

05463a1a8ab15818fb71d6bb4d3e5e4db4cb000a
