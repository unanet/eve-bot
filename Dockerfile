# Build ARGS
ARG BUILD_IMAGE
ARG GOOS="linux"
ARG GOARCH="amd64"
ARG VERSION=0.0.0
ARG GIT_BRANCH=""
ARG GIT_COMMIT_AUTHOR=""
ARG BUILD_HOST=""
ARG BUILDER=""
ARG GIT_COMMIT=""
ARG BUILD_DATE=""
ARG PRERELEASE=""

##########################################
# STEP 1 Unit Test & Build the Binary
##########################################
FROM ${BUILD_IMAGE} AS builder

# create appuser.
RUN adduser -D -g '' appuser

# set app working dir
WORKDIR $GOPATH/eve-bot

# Copy The source assets from the CWD (project root) into the container WORKDIR ($GOPATH/eve-bot)
COPY . .

# Unit Test 
# NOTE: (issues here with race due to alpine base gcc musl libs)
RUN GOOS=${GOOS} GOARCH=${GOARCH} go test -v ./...

# Golang buildtime ldflags
ENV LDFLAGS=" -X main.BuildHost=${BUILD_HOST} \
    -X main.GitBranch=${GIT_BRANCH} \
    -X main.Builder=${BUILDER} \
    -X main.Version=${VERSION} \
    -X main.BuildDate=${BUILD_DATE} \
    -X main.GitCommit=${GIT_COMMIT} \
    -X main.GitCommitAuthor=${GIT_COMMIT_AUTHOR} \
    -X main.VersionPrerelease=${PRERELEASE} "

# Build the binary
RUN GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" -o /go/bin/eve-bot ./cmd/eve-bot/


######################################
# STEP 2 build a smaller runtime image
######################################
FROM scratch

# Import assets from the build stage image
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/eve-bot /go/bin/eve-bot

# Use the unprivileged user (created in the build stage)
USER appuser

WORKDIR /go/bin

# Set the entrypoint to the golang executable binary
ENTRYPOINT ["/go/bin/eve-bot"]

# Expose the service ports (3000 for app and 3001 for metrics)
EXPOSE 3000
EXPOSE 3001

# Setup Container HealthCheck
HEALTHCHECK --interval=1m --timeout=2s --start-period=10s \
    CMD curl -f http://localhost:3000/ || exit 1
