# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS go-base

WORKDIR /src/server

COPY server/go.mod server/go.sum ./
RUN go mod download

COPY server/ ./

# Build both downloadable agents once, before selecting the image's architecture.
FROM go-base AS agent-build

RUN mkdir -p /out/agents \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/agents/netstamp-agent-linux-amd64 ./cmd/agent \
    && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o /out/agents/netstamp-agent-linux-arm64 ./cmd/agent

# Inherit the agents' Go build cache and cross-compile with the native toolchain.
FROM agent-build AS build

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/controller ./cmd/controller \
    && CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.22

RUN addgroup -S netstamp && adduser -S -G netstamp netstamp

WORKDIR /app

COPY --from=build /out/controller /app/controller
COPY --from=agent-build /out/agents /app/agents
COPY --from=build /out/migrate /app/migrate
COPY --from=build /src/server/db/migrations /app/db/migrations

USER netstamp:netstamp

EXPOSE 8080

ENTRYPOINT ["/app/controller"]
