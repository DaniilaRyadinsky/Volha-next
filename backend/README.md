# Volha Backend

Backend for the Volha furniture catalog.

## Current Services

- `gateway`: public HTTP API, Swagger, static `/images`, request ID, Redis cache adapter, gRPC client.
- `product-service`: gRPC catalog service, product/business logic, PostgreSQL repositories.
- `volha-migrator`: one-shot database migrator that runs before `product-service` in Compose.
- `dumper`: development backup helper.

## Local Checks

- `make format`
- `make build`
- `make vet`
- `make test`
- `make test-race`
- `make test-integration-compile`

Integration tests require PostgreSQL/Redis and live under `test/integration`. They are not run by the default `make test`.

`make up` builds the base Docker stack, waits for PostgreSQL health, runs migrations through `volha-migrator`, and only then starts `volha-product-service` and `volha-gateway`.

## Observability

Stage 3 adds technical observability:

- metrics: Go apps -> Prometheus -> Grafana
- logs: `slog` JSON -> stdout -> Grafana Alloy -> Loki -> Grafana
- traces: OpenTelemetry OTLP/gRPC -> Alloy -> Tempo -> Grafana
- alerts: Prometheus rules -> Alertmanager

The base Compose stack runs with `OTEL_ENABLED=false` and does not require Alloy. Use `make config-observability` to validate the combined Compose config and `make up-observability` to run the local stack. See `docs/observability.md`.

## Stage 2 Breaking Changes

- Removed legacy `/api/product/search`.
- Removed `views` and `is_favorite` from product API/schema.
- Removed `/api/auth/check` and `admin_pw` cookie/query auth.
- Removed generic `/api/message` mail endpoint.
- Product list pagination now uses `offset` and `limit`.

Future search, auth, favorites, product views, orders, notifications, and analytics will be implemented in separate stages.
