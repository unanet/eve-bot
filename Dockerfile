FROM unanet-docker.jfrog.io/alpine-base
RUN apk --no-cache add ca-certificates
ADD ./bin/eve-bot /app/eve-bot
WORKDIR /app
CMD ["/app/eve-bot"]
