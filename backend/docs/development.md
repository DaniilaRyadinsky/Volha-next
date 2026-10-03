# Development

## CURRENT

Use:

- `make fmt`
- `make build`
- `make vet`
- `make test-unit`
- `make test-unit-race`

Configuration is loaded from YAML plus environment overrides. Secrets must come from environment variables or an untracked `.env` file.

Redis is a cache dependency for the gateway. A Redis startup or runtime outage should degrade cache reads/writes only; catalog requests must still fall back to product-service.

## DATABASE RESET AFTER STAGE 2

Stage 2 squashed all development migrations into one initial schema:

- `build/volha/migrations/000001_init.up.sql`
- `build/volha/migrations/000001_init.down.sql`

Reason: the project has no production database, and previous migration history was development-only.

After Stage 2, local development DBs should be deleted and recreated, or the Docker volume should be removed. Do not try to migration-upgrade an old local schema.

Migration commands resolve `build/volha/migrations` from the repository root by default. Set `MIGRATIONS_PATH` only when running the migrator from a packaged/container path.

## FUTURE MIGRATION POLICY

Stage 2 is the last point where migration squash is allowed by default.

From now on, treat `000001_init` as the baseline. Stage 3+ schema changes should add new migrations:

- `000002_...`
- `000003_...`

Do not edit the initial migration for new features unless there is a separate explicit decision to reset the baseline again.

## TEMPORARY DEVELOPMENT ADMIN AUTH

Admin write endpoints use temporary development-only protection:

`Authorization: Bearer <DEV_ADMIN_TOKEN>`

This works only when `APP_ENV` is `local`, `development`, or `DEV`. Production mode rejects this mechanism.

This is NOT the production security model. It exists only until the Auth stage.
