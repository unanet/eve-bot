# eve-bot

![pipeline status](https://gitlab.unanet.io/devops/eve-bot/badges/master/pipeline.svg)

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
