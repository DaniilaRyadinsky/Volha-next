# Volha

Furniture catalog with a Go backend and frontend applications in a pnpm workspace.

## Repository Layout

```text
backend/          Go services, protobuf, service configs, migrations and Go tests
frontend/         Frontend applications and shared packages
.github/          GitHub Actions
ci/               Jenkins CI configuration
deploy/           Application deployment (reserved)
infrastructure/   Platform services; monitoring lives in platform/monitoring
ops/              Operational data; local database dumps live in database/dumps
tests/            Shared end-to-end and load tests (reserved)
docs/             Project documentation
compose.yaml      Local application stack
Makefile          Commands run from the repository root
```

Reserved directories are left empty without nested placeholders. Git does not track empty directories.

Backend-specific files stay in `backend`: `api` contains protobuf contracts and generated Go code, `configs` contains service YAML files, `docs` contains the generated Swagger Go package and API schemas, and `test/integration` contains Go integration tests. Unit tests live next to their packages.

The frontend workspace retains its existing `apps/web`, `apps/admin`, and `packages` structure.

## Local Setup

Run commands from the repository root. Copy `.env.example` to `.env` and set local credentials before starting services. The shared `.env` is ignored by Git and is also loaded when Go services run from `backend`.

```bash
cp .env.example .env
make config
make up
```

`backend/Dockerfile` has build targets `gateway` (the default), `auth-service`, `product-service`, `migrator`, and `dumper`. Compose selects the matching target for each service. The Compose project name remains `volha` so existing named volumes keep the same names.

## Current Services

- `gateway`: public HTTP API, Swagger, static `/images`, request ID, Redis cache adapter, gRPC client.
- `product-service`: gRPC catalog service, product/business logic, PostgreSQL repositories.
- `auth-service`: gRPC authentication service.
- `volha-migrator`: one-shot database migrator that runs before `product-service` in Compose.
- `dumper`: development backup helper.

## Local Checks

- `make format`
- `make build`
- `make vet`
- `make test`
- `make test-race`
- `make test-integration-compile`
- `make test-integration` (PostgreSQL + Redis)
- `make test-e2e` / `make test-regression` (PostgreSQL)

Integration, E2E and regression tests live under `backend/test/integration`; setup and coverage are described in [their README](backend/test/integration/README.md). They are not run by the default `make test`.

The GitHub Actions and local Jenkins quality gates are described in [docs/CI-CD.md](docs/CI-CD.md). Integration tests are planned for a separate CDL pipeline.

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

Future search, favorites, product views, orders, notifications, and analytics will be implemented in separate stages.
