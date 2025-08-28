FROM golang:1.24.0-alpine AS builder

RUN apk add --no-cache git && \
    go install github.com/TBXark/mcp-proxy@latest && \
    go install github.com/mark3labs/mcp-filesystem-server@latest

FROM alpine:3.19

RUN apk --no-cache add ca-certificates && \
    addgroup -S mcp && \
    adduser -S -G mcp mcp && \
    mkdir -p /app /data && \
    chown -R mcp:mcp /app /data

# Set working directory
WORKDIR /app

COPY --from=builder /go/bin/mcp-proxy /usr/local/bin/
COPY --from=builder /go/bin/mcp-filesystem-server /usr/local/bin/

COPY config.json /app/config.json

USER mcp

ENTRYPOINT [ "mcp-proxy", "/app/config.json" ]
