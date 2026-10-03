# Security

## CURRENT

Secrets are not kept in tracked YAML configs. `.env` is ignored; `.env.example` contains placeholders only.

Gateway is the only publicly exposed application service in Docker Compose. PostgreSQL, Redis, and product-service gRPC stay internal by default.

## TEMPORARY DEVELOPMENT ADMIN AUTH

Administrative write endpoints are protected only in development by:

`Authorization: Bearer <DEV_ADMIN_TOKEN>`

Rules:

- accepted only for `APP_ENV=local`, `APP_ENV=development`, or `APP_ENV=DEV`;
- rejected in production mode;
- token is not logged;
- token is not stored in cookies;
- there is no login endpoint or session.

This is NOT production authentication. It replaces the removed `admin_pw` query/cookie flow until auth-service is implemented.

## REMOVED LEGACY SECURITY RISKS

- Removed `/api/auth/check?admin_pw=...`.
- Removed raw password cookie.
- Removed generic MailerSend endpoint and related API key usage.

Any credentials that were previously committed should be rotated.

## FUTURE

Planned areas: auth-service, secrets management, TLS, reverse proxy hardening, rate limiting, admin security rewrite, container scanning, dependency scanning, SAST, backup encryption, deployment rollback, smoke testing, and load testing.
