# Volha Backend Architecture

## Current Core

Gateway is the only public application entrypoint. It owns HTTP transport, `/api`, `/images`, Swagger, health endpoints, request IDs, Redis cache access, and product-service gRPC client calls.

Product Service owns catalog business logic, gRPC transport, service validation, transaction ownership, typed repositories, and PostgreSQL access.

PostgreSQL is the durable product/catalog/content store.

Redis is an optional dictionary cache. Redis outage degrades cache behavior but should not make the catalog unavailable if product-service is healthy.

```text
Client
  -> Gateway
    -> Product Service
      -> PostgreSQL
    -> Redis cache
```

## Observability

Stage 3 adds technical observability:

- Prometheus
- Grafana
- Loki
- Tempo
- Alertmanager
- Grafana Alloy
- postgres-exporter
- redis-exporter
- blackbox-exporter
- cAdvisor
- optional node-exporter

```text
Gateway/Product Service
  -> /metrics -> Prometheus -> Grafana
  -> slog stdout -> Alloy -> Loki -> Grafana
  -> OTLP/gRPC -> Alloy -> Tempo -> Grafana
  -> Prometheus rules -> Alertmanager
```

Application metrics, logs, and traces are separate pipelines. Business analytics is out of scope.

## Data Flow

```text
HTTP handler
  -> gateway product gRPC client
  -> product-service gRPC transport
  -> service
  -> typed repository interface/provider
  -> PostgreSQL repository
```

SQL belongs in repository code only.

## Stage 3 Baseline

Legacy catalog search remains removed. Smart Search will be designed later as a separate stage.

Product views remain removed. Future product views should be emitted as business events through a future analytics-service.

Favorites remain removed. Favorites need users/auth and a user-product relation.

The old generic mail endpoint remains removed. Future order/callback notifications should flow through order-service, outbox, and notification-service.

## Future Architecture

These services are roadmap placeholders only:

- Auth Service
- Order Service
- Search Service
- Content Service
- Analytics Service
- Notification Service

No CI/CD, Kubernetes, business analytics, orders, auth, smart search, or user event tracking is implemented in Stage 3.
