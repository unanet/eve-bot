# Build ARGS
ARG GOOS="linux"
ARG GOARCH="amd64"
ARG GO_VERSION="1.14.1-alpine"
ARG VERSION=0.0.0
ARG GIT_BRANCH=""
ARG GIT_COMMIT_AUTHOR=""
ARG BUILD_HOST=""
ARG BUILDER=""
ARG GIT_COMMIT=""
ARG BUILD_DATE=""
ARG PRERELEASE=""


##########################################
# STEP 1 build binary in Build Stage Image
##########################################
FROM golang:${GO_VERSION} AS builder

LABEL maintainer="Unanet DevOps <ops@unanet.io>"

# Golang buildtime ldflags
ENV LDFLAGS=" -X main.BuildHost=${BUILD_HOST} \
    -X main.GitBranch=${GIT_BRANCH} \
    -X main.Builder=${BUILDER} \
    -X main.Version=${VERSION} \
    -X main.BuildDate=${BUILD_DATE} \
    -X main.GitCommit=${GIT_COMMIT} \
    -X main.GitCommitAuthor=${GIT_COMMIT_AUTHOR} \
    -X main.VersionPrerelease=${PRERELEASE} "

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
# tzdata is for timezone data
RUN apk update && apk add --no-cache gcc musl-dev git ca-certificates tzdata && update-ca-certificates

# create appuser.
RUN adduser -D -g '' appuser

# set app working dir
WORKDIR $GOPATH/eve-bot

# Copy The source assets from the CWD (project root) into the container WORKDIR ($GOPATH/eve-bot)
COPY . .

# Test It
RUN GOOS=${GOOS} GOARCH=${GOARCH} go test -v ./...

# Build It
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" -o /go/bin/eve-bot ./cmd/eve-bot/


######################################
# STEP 2 build a smaller runtime image
######################################
FROM scratch

LABEL maintainer="Unanet DevOps <ops@unanet.io>"

# Import assets from the build stage image
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/eve-bot /go/bin/eve-bot

# Use the unprivileged user (created in the build stage image)
USER appuser

WORKDIR /go/bin

# Set the entrypoint to the golang executable binary
ENTRYPOINT ["/go/bin/eve-bot"]

# Expose the service ports (3000 for app and 3001 for metrics)
EXPOSE 3000
EXPOSE 3001

# Setup Container HealthCheck
HEALTHCHECK --interval=1m --timeout=2s --start-period=10s \
    CMD curl -f http://localhost:3000/ping || exit 1
