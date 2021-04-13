FROM unanet-docker.jfrog.io/alpine-base

ENV SERVICE_NAME eve-bot
ENV LOG_LEVEL debug
ENV EVEBOT_PORT 3000
ENV EVEBOT_METRICS_PORT 3001
ENV EVEBOT_EVEAPI_BASE_URL http://eve-api-v1:3000
ENV EVEBOT_EVEAPI_CALLBACK_URL http://eve-bot-v1:3000/eve-callback
ENV EVEBOT_SLACK_CHANNELS_AUTH my-evebot,evebot-tests,hydra,admin-ci
ENV EVEBOT_SLACK_CHANNELS_MAINTENANCE my-evebot,evebot-tests
ENV EVEBOT_SLACK_AUTH_ENABLED true
ENV VAULT_ADDR https://vault.unanet.io
ENV VAULT_ROLE k8s-devops
ENV VAULT_K8S_MOUNT kubernetes

ADD ./bin/eve-bot /app/eve-bot
WORKDIR /app
CMD ["/app/eve-bot"]

HEALTHCHECK --interval=1m --timeout=2s --start-period=10s \
    CMD curl -f http://localhost:${EVE_METRICS_PORT}/ || exit 1
