# CLAUDE.md

## Tooling

`brew install` mockery, sqlc, direnv, and Docker.

## Local Postgres

`make db-up` brings up `postgres:16` configured to match what `pkg/testdb` expects (OS-user superuser, password `postgres`, port from `PGPORT` — defaults to `54329` to avoid clashes; `make test` and the committed `.envrc` both wire it through).

## Architecture: DDD-lite

**Imports flow inward only.** Allowed imports between layers under `pkg/withoutings/`:

- `domain/<aggregate>` → nothing in `pkg/withoutings/`
- `app/{command,query,service}` → `domain/*`
- `app/app.go` → `adapter/*`, `domain/*` (sole composition root)
- `adapter/<aggregate>` → `domain/*` (the adapter implements `domain.Repo`)
- `port/`, `pkg/worker/` → `app`, plus `adapter/topic` for shared topic-name constants

Gotchas:

- **`App.Validate()` is called from `web.Router`** and panics on any nil wire-up — adding a new command/query means adding it to both `Commands{}`/`Queries{}` AND `Validate()`.
- **The `withings` adapter is an HTTP client**, not a Postgres adapter — it implements `domain/withings.Repo` and does not touch `pkg/db`.
- **Two DB handles, on purpose**: `app.NewApplication` constructs both a `*pgxpool.Pool` (for everything else) and a `*sql.DB` (required by `watermill-sql/v3`). Don't unify them.
- **`WithTx(pgx.Tx)` lives on the adapter struct, not on the domain `Repo` interface.** Multi-aggregate writes should run inside a single pgx transaction by composing `WithTx`-wrapped adapters.

## Pointers

- `pkg/integrationtest` brings up a fresh DB + fully-wired `*app.App` + worker + router — use it for cross-layer tests rather than reassembling pieces.
