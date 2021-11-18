<div id="top"></div>

# eve-bot

![pipeline status](https://github.com/unanet/eve-bot/badges/master/pipeline.svg)

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
       <a href="#about-the-project">About The Project</a>
    </li>
    <li><a href="#building">Building</a></li>
    <li><a href="#running">Running</a></li>
    <li><a href="#environment-variables">Environment Variables</a></li>
    <li>
      <a href="#configuration">Getting Started</a>
      <ul>
        <li><a href="#slack">Slack</a></li>
      </ul>
    </li>
  </ol>
</details>



## About The Project

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
EVEBOT_EVEAPI_BASE_URL=""
EVEBOT_EVEAPI_CALLBACK_URL=""
EVEBOT_EVEAPI_ADMIN_TOKEN=""
EVEBOT_IDENTITY_CONN_URL=""
EVEBOT_IDENTITY_CLIENT_ID=""
EVEBOT_IDENTITY_CLIENT_SECRET=""
EVEBOT_IDENTITY_REDIRECT_URL=""
EVEBOT_AWS_REGION=""
EVEBOT_LOGGING_DASHBOARD_BASE_URL=""
EVEBOT_USER_TABLE_NAME=""
EVEBOT_DEVOPS_MONITORING_CHANNEL=""
```

## Getting Started

### Slack

#### Slack Environment Variables

```bash
EVEBOT_SLACK_SIGNING_SECRET=""
EVEBOT_SLACK_VERIFICATION_TOKEN=""
EVEBOT_SLACK_OAUTH_ACCESS_TOKEN=""
```

* [Create an App ( From Scratch )](https://api.slack.com/apps)
* App summary will have the `Signing Secret` which will be `EVEBOT_SLACK_SIGNING_SECRET` and the `Verification Token` for `EVEBOT_SLACK_VERIFICATION_TOKEN`
* Enable Incoming web hooks
* Add bot via OAuth & Permissions to your channel
    * Copy Bot User OAuth Token, this will be the value for `EVEBOT_SLACK_OAUTH_ACCESS_TOKEN`


```yaml
_metadata:
  major_version: 1
  minor_version: 1
display_information:
  name: {{bot-name}}
features:
  bot_user:
    display_name: {{bot-name}}
    always_online: true
oauth_config:
  scopes:
    user:
      - users:read
    bot:
      - app_mentions:read
      - incoming-webhook
      - users:read
settings:
  event_subscriptions:
    request_url: https://{{domain}}/slack-events
    bot_events:
      - app_mention
  org_deploy_enabled: false
  socket_mode_enabled: false
  token_rotation_enabled: false
```
