# 0002: Use goose for database migrations

Status: Accepted
Date: 2026-11-15

## Context

The database structure lives in SQL files in `migrations/`. At first, Docker ran them automatically from `/docker-entrypoint-initdb.d`, but only when the database was created empty. Later files, like the tasks table, would never run on an existing database, and I had to apply them by hand. I also had no record of which files were already applied.

## Decision

Use goose. Each migration is a `.sql` file with `-- +goose Up` and `-- +goose Down` sections. goose keeps the applied versions in a `goose_db_version` table and runs only the new ones. It is a separate command-line program, not part of the Sokin code. I run it from the project root with the `GOOSE_*` variables built from `.env`.

The Docker `initdb` mount is removed, so there is one way to apply migrations. Migrations 001 and 002 were applied by hand before goose, so I marked them as applied in `goose_db_version`.

## Alternatives considered

- **Keep `initdb` only:** nothing to install, but it never updates an existing database.
- **Run SQL files by hand with `psql`:** simple, but easy to forget a file or run one twice.
- **golang-migrate:** a very similar tool. I picked goose because its SQL files are simple and it works well as a CLI.
- **An ORM with auto-migrate:** hides the SQL, and I want to learn and review the SQL myself.

## Consequences

**Good**
- Every change to the structure is a numbered file in Git, and the database knows which ones it has.
- The same command will work in the cloud and in CI in Phase 8.
- Each migration can be undone with `goose down`.

**Bad**
- goose must be installed separately (`go install`), and the connection variables must be set before using it.
- `goose down` drops columns and tables together with their data. It is only safe before real data exists, so I take a backup before every migration.