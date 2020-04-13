FROM golang:1.14.2-alpine AS builder
WORKDIR /src/
COPY cmd /src/cmd
COPY internal /src/internal
COPY pkg /src/pkg
COPY go.mod /src/
COPY go.sum /src/
RUN --mount=type=ssh go build -o eve-bot ./cmd/eve-bot/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /src/eve-bot /app/eve-bot
CMD ["/app/eve-bot"]
#COPY app.go    .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#COPY --from=builder /go/src/github.com/alexellis/href-counter/app .
#CMD ["./app"]
#
#
#######################################
## STEP 2 build a smaller runtime image
#######################################
#FROM scratch
#
## Import assets from the build stage image
#COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY --from=builder /etc/passwd /etc/passwd
#COPY --from=builder /go/bin/eve-bot /go/bin/eve-bot
#
#USER appuser
#
#WORKDIR /go/bin
#ENTRYPOINT ["/go/bin/eve-bot"]
#EXPOSE 3000
#EXPOSE 3001
#
#HEALTHCHECK --interval=1m --timeout=2s --start-period=10s \
#    CMD curl -f http://localhost:3000/ || exit 1
