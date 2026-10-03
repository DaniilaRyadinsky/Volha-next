# Deployment

## CURRENT

Docker Compose builds and runs gateway, product-service, PostgreSQL, Redis, and dumper on an internal network. Only gateway publishes a host port by default.

Dockerfiles use multi-stage builds, non-root runtime users for Go services, small runtime images, and explicit image tags instead of `latest`.

Uploaded application images are persisted in the Docker named volume `volha-images` and served by the gateway under `/images`.

## FUTURE

Roadmap areas: CI/CD, staging, production, reverse proxy, TLS, healthchecks, Prometheus, Grafana, Loki, Tempo, Alertmanager, Grafana Alloy, OpenTelemetry, exporters, centralized logging, distributed tracing, alerting, backup, restore, security scanning, dependency scanning, container scanning, deployment rollback, load testing, and smoke testing.

No CI/CD pipeline, Kubernetes, reverse proxy, or TLS configuration is implemented yet.
