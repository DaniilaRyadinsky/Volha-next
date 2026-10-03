FROM golang:1.25.4-alpine3.22 AS builder

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Version=${VERSION} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Commit=${COMMIT} -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.BuildTime=${BUILD_TIME}" \
  -o /app/product-service ./cmd/product-service/app/main.go

FROM alpine:3.22
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
LABEL org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="volha-product-service"
WORKDIR /app
RUN addgroup -S volha && adduser -S volha -G volha && mkdir -p /app/configs && chown -R volha:volha /app
COPY --from=builder /app/product-service /app/product-service

VOLUME /app/configs
USER volha

EXPOSE 8008

ENTRYPOINT ["/app/product-service"]
