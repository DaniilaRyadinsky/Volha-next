# Testing

## Unit

Unit tests live next to internal packages and do not require external dependencies.

Run:

- `make test-unit`
- `make test-unit-race`
- `make test-cover`

## Integration

Integration tests live in `test/integration` and require PostgreSQL and/or Redis.

Integration config reads credentials from environment variables. Prefer `TEST_POSTGRES_*` and `TEST_REDIS_*` for isolated local runs; otherwise the helpers fall back to the normal `POSTGRES_*` and `REDIS_*` variables.

Run explicitly only when the integration environment is available:

- `make test-integration`
- `make test-integration-compile` compiles the integration suite without running DB/Redis workloads.

Existing DB/Redis tests were moved here and are guarded by the `integration` build tag.

## E2E

E2E is future work. It may use Postman/Newman later, but no E2E suite is implemented now.

## All Tests

`make test-all` runs unit and integration targets. Do not use it unless PostgreSQL/Redis are available.
