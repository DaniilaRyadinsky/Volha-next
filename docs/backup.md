# Backup And Restore

## CURRENT

The existing dumper command is preserved. Runtime dumps are stored in `ops/database/dumps` and ignored by git.

The dumper reads database settings from the process environment; a physical `.env` file is optional. Dumps are written to a temporary `.tmp` path and renamed only after `pg_dump` succeeds, so failed backups do not leave valid-looking partial dump files.

## FUTURE

Planned areas: scheduled backups, retention policy, off-host storage, restore testing, backup monitoring, encryption, and documented restore drills.

No remote backup storage or restore automation is implemented yet.
