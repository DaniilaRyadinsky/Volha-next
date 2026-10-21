# Build a service with --target gateway, auth-service, product-service, migrator or dumper.

FROM golang:1.25.4-alpine3.22 AS build-base

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM build-base AS auth-service-builder
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Version=${VERSION} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Commit=${COMMIT} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.BuildTime=${BUILD_TIME}" \
  -o /app/auth-service ./cmd/auth-service

FROM alpine:3.22 AS auth-service

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

LABEL org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="volha-auth-service"

WORKDIR /app
RUN addgroup -S volha && adduser -S volha -G volha && mkdir -p /app/configs && chown -R volha:volha /app
COPY --from=auth-service-builder /app/auth-service /app/auth-service

VOLUME /app/configs
USER volha

EXPOSE 8009 9093

ENTRYPOINT ["/app/auth-service"]

FROM build-base AS product-service-builder
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Version=${VERSION} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Commit=${COMMIT} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.BuildTime=${BUILD_TIME}" \
  -o /app/product-service ./cmd/product-service/app/main.go

FROM alpine:3.22 AS product-service
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
LABEL org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="volha-product-service"
WORKDIR /app
RUN addgroup -S volha && adduser -S volha -G volha && mkdir -p /app/configs && chown -R volha:volha /app
COPY --from=product-service-builder /app/product-service /app/product-service

VOLUME /app/configs
USER volha

EXPOSE 8008

ENTRYPOINT ["/app/product-service"]

FROM build-base AS migrator-builder
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Version=${VERSION} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Commit=${COMMIT} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.BuildTime=${BUILD_TIME}" \
  -o /app/volha-migrator ./cmd/product-service/migrator

FROM alpine:3.22 AS migrator

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

LABEL org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="volha-migrator"

WORKDIR /app
RUN addgroup -S volha && adduser -S volha -G volha && mkdir -p /migrations && chown -R volha:volha /app /migrations
COPY --from=migrator-builder /app/volha-migrator /app/volha-migrator
COPY migrations /migrations

ENV MIGRATIONS_PATH=/migrations
USER volha

ENTRYPOINT ["/app/volha-migrator"]

FROM build-base AS dumper-builder
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Version=${VERSION} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Commit=${COMMIT} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.BuildTime=${BUILD_TIME}" \
  -o /app/dumper ./cmd/dumper

FROM postgres:18-alpine AS dumper
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
LABEL org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="volha-dumper"

WORKDIR /app


RUN mkdir -p /dumps

COPY --from=dumper-builder /app/dumper /app/dumper

ENTRYPOINT ["/app/dumper"]

FROM build-base AS gateway-builder
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Version=${VERSION} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Commit=${COMMIT} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.BuildTime=${BUILD_TIME}" \
  -o /app/volha-gateway ./cmd/gateway/main.go

FROM alpine:3.22 AS gateway
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
LABEL org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="volha-gateway"
WORKDIR /app
RUN addgroup -S volha && adduser -S volha -G volha && mkdir -p /app/configs /app/images && chown -R volha:volha /app
COPY --from=gateway-builder /app/volha-gateway /app/volha-gateway

VOLUME /app/configs
VOLUME /app/images
USER volha

EXPOSE 8080

ENTRYPOINT ["/app/volha-gateway"]
