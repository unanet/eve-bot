FROM unanet-docker.jfrog.io/alpine-base

ENV EVE_PORT 8080
ENV EVE_METRICS_PORT 3000
ENV EVE_SERVICE_NAME eve-bot

ADD ./bin/eve-bot /app/eve-bot
WORKDIR /app
CMD ["/app/eve-bot"]

HEALTHCHECK --interval=1m --timeout=2s --start-period=10s \
    CMD curl -f http://localhost:${EVE_METRICS_PORT}/ || exit 1
