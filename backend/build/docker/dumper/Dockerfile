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
  -o /app/dumper ./cmd/dumper


FROM postgres:18-alpine
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
LABEL org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.title="volha-dumper"

WORKDIR /app


RUN mkdir -p /dumps

COPY --from=builder /app/dumper /app/dumper

ENTRYPOINT ["/app/dumper"]
