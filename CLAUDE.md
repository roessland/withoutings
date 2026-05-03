# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Tooling

`brew install` mockery, sqlc, direnv, and Docker. `.mockery.yaml` is v3 schema.

## Local Postgres

`make db-up` brings up `postgres:16` configured to match what `pkg/testdb` expects (OS-user superuser, password `postgres`, port from `PGPORT` — defaults to `54329` to avoid clashes; `make test` and the committed `.envrc` both wire it through).

## Architecture: DDD-lite (threedots.tech)

**Imports flow inward only.** This is the rule that breaks most often when changes get made carelessly:

```
port  ─┐
       ├──►  app  ──►  domain
adapter┘            ▲
                    │
              (adapter implements domain.Repo interfaces)
```

Layers under `pkg/withoutings/`: `domain/<aggregate>`, `app/{command,query,service}`, `adapter/<aggregate>`, `port/`. Aggregates are `account`, `subscription`, `withings`.

Subtleties that are easy to violate:

- **`app/app.go` is the sole composition root** that imports `adapter/*`. Command/query/service code must depend on `domain.Repo` interfaces, never on adapter types.
- **`App.Validate()` is called from `web.Router`** and panics on any nil wire-up — adding a new command/query means adding it to both `Commands{}`/`Queries{}` AND `Validate()`.
- **The `withings` adapter is an HTTP client**, not a Postgres adapter — it implements `domain/withings.Repo` and does not touch `pkg/db`.
- **`adapter/topic/` is the one legitimate exception** to the inward-only rule: `port/` and `worker/` may import it for shared topic-name constants.
- **Two DB handles, on purpose**: `app.NewApplication` constructs both a `*pgxpool.Pool` (for everything else) and a `*sql.DB` (required by `watermill-sql/v3`). Don't unify them.
- **`WithTx(pgx.Tx)` lives on the adapter struct, not on the domain `Repo` interface.** Multi-aggregate writes should run inside a single pgx transaction by composing `WithTx`-wrapped adapters.
- **Domain event payloads stay minimal (UUID references)** until the TODO in `domain/subscription/subscription_events.go` separates domain events from watermill/JSON encoding.

### Adding a feature

1. Domain types and `Repo` interface in `domain/<aggregate>/`.
2. Adapter impl in `adapter/<aggregate>/` + SQL in `pkg/db/<aggregate>.sql` → `make generate-sqlc`.
3. Command or query handler in `app/{command,query}/` depending on the domain interface; wire into `app.Commands`/`app.Queries` AND `Validate()`.
4. HTTP handler in `port/`, route in `pkg/web/routes.go`.
5. Async work: publish topic from `adapter/topic/topics.go`, consume in `pkg/worker/worker.go`.

## Pointers

- `pkg/integrationtest` brings up a fresh DB + fully-wired `*app.App` + worker + router — use it for cross-layer tests rather than reassembling pieces.
