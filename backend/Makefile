.DEFAULT_GOAL := help

VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_TIME ?= unknown
export VERSION
export COMMIT
export BUILD_TIME

LDFLAGS := -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Version=$(VERSION) -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.Commit=$(COMMIT) -X github.com/autumnterror/volha-backend/internal/platform/buildinfo.BuildTime=$(BUILD_TIME)

COMPOSE := docker compose -f ./build/volha/docker-compose.yml
COMPOSE_OBSERVABILITY := docker compose -f ./build/volha/docker-compose.yml -f ./deployments/observability/docker-compose.yml
COMPOSE_OBSERVABILITY_LINUX := docker compose -f ./build/volha/docker-compose.yml -f ./deployments/observability/docker-compose.yml -f ./deployments/observability/docker-compose.linux.yml --profile linux-host

.PHONY: help format vet build run-gateway run-product-service generate generate-proto generate-swagger test test-race test-cover test-integration test-integration-compile test-all docker-build up down logs ps config up-observability down-observability logs-observability ps-observability config-observability up-observability-linux config-observability-linux migrate-up migrate-down

help:
	@echo "Volha backend targets:"
	@echo "  format                         go fmt ./..."
	@echo "  vet                            go vet ./..."
	@echo "  build                          build all Go packages"
	@echo "  run-gateway                    run gateway locally"
	@echo "  run-product-service            run product service locally"
	@echo "  generate                       regenerate proto and swagger"
	@echo "  generate-proto                 regenerate protobuf Go code"
	@echo "  generate-swagger               regenerate Swagger docs"
	@echo "  test                           run Go tests"
	@echo "  test-race                      run Go tests with race detector"
	@echo "  test-cover                     run Go tests with coverage"
	@echo "  test-integration               run integration tests"
	@echo "  test-integration-compile       compile integration tests only"
	@echo "  test-all                       run unit and integration tests"
	@echo "  docker-build                   build Docker images"
	@echo "  up/down/logs/ps/config         manage base Compose stack"
	@echo "  up-observability/...           manage Compose stack with observability"
	@echo "  migrate-up/migrate-down        run DB migrations through Docker"

format:
	go fmt ./...

vet:
	go vet ./...

build:
	go build -ldflags "$(LDFLAGS)" ./...

run-gateway:
	go run ./cmd/gateway -config ./build/volha/configs/volha-gateway.yaml

run-product-service:
	go run ./cmd/product-service/app -config ./build/volha/configs/product-service.yaml

generate: generate-proto generate-swagger

generate-proto:
	docker run --rm -v "$(CURDIR):/workspace" -w /workspace bufbuild/buf:1.50.0 generate api/proto/filez --template '{"version":"v2","plugins":[{"remote":"buf.build/protocolbuffers/go:v1.36.10","out":"api/proto/gen","opt":["paths=source_relative"]},{"remote":"buf.build/grpc/go:v1.5.1","out":"api/proto/gen","opt":["paths=source_relative"]}]}'

generate-swagger:
	swag init --dir ./cmd/gateway,./internal/gateway/net/handlers,./internal/gateway/net/view,./internal/gateway/model,./api/proto/gen --output ./docs

test:
	go test ./...

test-race:
	go test -race ./...

test-cover:
	go test -cover ./...

test-integration:
	go test -tags=integration ./test/integration/...

test-integration-compile:
	go test -tags=integration ./test/integration/... -run '^$$'

test-all: test test-integration

docker-build:
	$(COMPOSE) build

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

config:
	$(COMPOSE) config

up-observability:
	$(COMPOSE_OBSERVABILITY) up -d --build

down-observability:
	$(COMPOSE_OBSERVABILITY) down

logs-observability:
	$(COMPOSE_OBSERVABILITY) logs -f

ps-observability:
	$(COMPOSE_OBSERVABILITY) ps

config-observability:
	$(COMPOSE_OBSERVABILITY) config

up-observability-linux:
	$(COMPOSE_OBSERVABILITY_LINUX) up -d --build

config-observability-linux:
	$(COMPOSE_OBSERVABILITY_LINUX) config

migrate-up:
	$(COMPOSE) run --rm volha-migrator -type up

migrate-down:
	$(COMPOSE) run --rm volha-migrator -type down
