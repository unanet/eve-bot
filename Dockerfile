##########################################
# STEP 1 Setup the base image requirements
##########################################
ARG BUILD_IMAGE
FROM ${BUILD_IMAGE} AS builder

# create appuser.
RUN adduser -D -g '' appuser

# Add the binary
ADD build/eve-bot /go/bin/eve-bot


######################################
# STEP 2 build a smaller runtime image
######################################
FROM scratch

# Import assets from the build stage image
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/eve-bot /go/bin/eve-bot

USER appuser

WORKDIR /go/bin
ENTRYPOINT ["/go/bin/eve-bot"]
EXPOSE 3000
EXPOSE 3001

HEALTHCHECK --interval=1m --timeout=2s --start-period=10s \
    CMD curl -f http://localhost:3000/ || exit 1
