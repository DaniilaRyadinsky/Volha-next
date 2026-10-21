# Observability

Volha Stage 3 adds technical observability only. It is not business analytics: no funnels, product-view events, Yandex Metrika, order statistics, or recommendations.

## Architecture

Metrics:

```text
Go apps -> /metrics -> Prometheus -> Grafana
```

Logs:

```text
Go slog JSON -> stdout/stderr -> Grafana Alloy -> Loki -> Grafana
```

Traces:

```text
Go OpenTelemetry SDK -> OTLP/gRPC -> Grafana Alloy -> Tempo -> Grafana
```

Alerts:

```text
Prometheus rules -> Alertmanager -> Grafana / Alertmanager UI
```

Application logs do not use the OpenTelemetry Logs SDK. Logs stay structured `log/slog` records and are collected from container stdout by Alloy.

## Run Locally

Create `.env` from the root `.env.example` and set real local values. Grafana must use non-empty local credentials.
The base compose stack keeps tracing disabled (`OTEL_ENABLED=false`) and does not require Alloy. The observability overlay enables OTLP tracing for `volha-gateway` and `volha-product-service`.

```bash
make up-observability
make config-observability
make down-observability
```

Use Docker Compose directly if you intentionally need to delete local monitoring volumes.

Linux host metrics are optional:

```bash
make up-observability-linux
make config-observability-linux
```

## Local URLs

| Service | URL |
| --- | --- |
| Grafana | http://127.0.0.1:3000 |
| Prometheus | http://127.0.0.1:9090 |
| Alertmanager | http://127.0.0.1:9093 |
| Alloy UI | http://127.0.0.1:12345 |

Loki, Tempo, exporters, and application metrics ports are internal Docker endpoints only.

## Datasources

Grafana provisioning creates Prometheus, Loki, Tempo, and Alertmanager datasources. Loki has a derived `trace_id` field for log-to-trace navigation. Tempo is configured for trace-to-logs lookup using `service.name` and `trace_id`. Tempo service graph is intentionally disabled because Tempo metrics-generator is not configured in this stack.

## Dashboards

Provisioned dashboards:

- `Volha / Overview`
- `Volha / Gateway`
- `Volha / Product Service`
- `Volha / Infrastructure`

They cover service health, HTTP/gRPC rates and latency, error rates, DB pool usage, Redis cache results, container CPU/memory/network, PostgreSQL exporter metrics, and Redis exporter metrics.

## Alerts

Prometheus alert rules include gateway/product-service down, PostgreSQL exporter down, Redis exporter down, gateway readiness blackbox failure, HTTP 5xx rate above 5 percent, HTTP p95 above 1 second, and gRPC error rate above 5 percent. Ratio and latency alerts require minimum request traffic to avoid idle-stack noise.

Recording rules create 5-minute HTTP/gRPC request and error rates.

## Application Metrics

Custom metrics:

- `volha_http_requests_total` counter labels: `service`, `method`, `route`, `status`.
- `volha_http_request_duration_seconds` histogram labels: `service`, `method`, `route`, `status`.
- `volha_http_inflight_requests` gauge labels: `service`.
- `volha_grpc_server_requests_total` counter labels: `service`, `grpc_service`, `grpc_method`, `grpc_code`.
- `volha_grpc_server_request_duration_seconds` histogram labels: `service`, `grpc_service`, `grpc_method`, `grpc_code`.
- `volha_grpc_client_requests_total` counter labels: `service`, `grpc_service`, `grpc_method`, `grpc_code`.
- `volha_grpc_client_request_duration_seconds` histogram labels: `service`, `grpc_service`, `grpc_method`, `grpc_code`.
- `volha_redis_cache_operations_total` counter labels: `service`, `cache`, `operation`, `result`.
- `volha_build_info` gauge labels: `service`, `version`, `commit`.
- `volha_db_pool_open_connections` gauge labels: `service`.
- `volha_db_pool_in_use_connections` gauge labels: `service`.
- `volha_db_pool_idle_connections` gauge labels: `service`.
- `volha_db_pool_wait_count_total` counter labels: `service`.
- `volha_db_pool_wait_duration_seconds_total` counter labels: `service`.

Go and process collectors are also enabled.

## Cardinality Rules

Never use these as metric labels:

- `user_id`
- `product_id`
- `request_id`
- `trace_id`
- raw URL
- search query
- email
- IP address
- error message

Metrics labels must stay low-cardinality and bounded. High-cardinality identifiers belong in logs and traces. HTTP metrics use Echo normalized routes such as `/api/product/:id`, not raw URLs or query strings.

## Logging

`backend/internal/platform/logger` uses `log/slog`.

Config:

- `LOG_LEVEL=debug|info|warn|error`
- `LOG_FORMAT=text|json`

Fields:

- Always: timestamp, level, message, service, environment, version.
- From context: request_id, trace_id, span_id.
- Operational: operation, component, duration_ms, grpc_method, http_method, http_route, http_status, error.

Do not log Authorization headers, cookies, tokens, PostgreSQL passwords, Redis passwords, request bodies, or raw query strings.

Canonical service names are `volha-gateway` and `volha-product-service` across logs, metrics, traces, Loki labels, and Grafana queries.
Alloy parses application slog JSON for `volha-gateway`, `volha-product-service`, and `volha-dumper`, promoting only `service`, `environment`, and `level` to Loki labels. `request_id`, `trace_id`, and `span_id` remain structured fields for querying/correlation, not labels.

## Tracing

Config:

- `OTEL_ENABLED=true|false` (`false` in base compose, `true` in the observability overlay)
- `OTEL_EXPORTER_OTLP_ENDPOINT=alloy:4317`
- `OTEL_TRACES_SAMPLER_RATIO=0.1`

Tracing uses official Echo OpenTelemetry middleware with W3C trace-context propagation. Tracing is optional. If Alloy, Tempo, Loki, Prometheus, or Grafana are down, the catalog application must continue to work.

## pprof

`PPROF_ENABLED=true` exposes Go pprof handlers on the internal metrics server. It is disabled by default and not published to the host.

## Exporters

- `postgres-exporter` uses local database credentials from environment. Future production deployments should use least-privilege monitoring credentials.
- `redis-exporter` uses Redis credentials from environment.
- `blackbox-exporter` probes gateway liveness/readiness.
- `cAdvisor` collects container metrics with Docker host mounts.
- `node-exporter` is optional via Compose profile `linux-host`; Docker Desktop on Windows/macOS is not a Linux host equivalent.

## Troubleshooting

- Prometheus targets: `http://127.0.0.1:9090/targets`.
- Loki logs: Grafana Explore with datasource `Loki`.
- Tempo traces: Grafana Explore with datasource `Tempo`.
- Alloy status: `http://127.0.0.1:12345`.
- Empty dashboards: check app metrics endpoints `volha-gateway:9091` and `volha-product-service:9092`.
- Missing traces: check `OTEL_ENABLED=true` and Alloy OTLP receiver `alloy:4317`.
