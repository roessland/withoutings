# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common commands

```bash
make db-up                # start a local Postgres in Docker (idempotent)
make test                 # PGPORT=$(DB_PORT) go test -v -race -cover ./...
make build                # builds binaries via `go build ./...`
make generate-all         # `go generate ./...` (mockery + sqlc)
make generate-sqlc        # regenerate pkg/db from SQL after editing migrations or *.sql
make migrate              # source env.dev.sh && go run cmd/withoutings/*.go migrate
make run-dev              # generate-all + migrate + run server with env.dev.sh
make db-down              # stop the Postgres container, keep the volume
make db-destroy           # db-down + drop the volume

# Single tests
go test -run TestRandomNonceIsURLSafe ./pkg/withoutings/domain/withings/...    # pure unit, no DB
go test -race -run TestCreateOrUpdateAccountHandler ./pkg/withoutings/app/command/...  # needs db-up
```

### Running tests locally

`pkg/testdb` connects to Postgres with the **OS user as a superuser** and password `postgres` (see `pkg/testdb/testdb.go`). The connection string omits the port, so libpq/pgx falls back to `PGPORT`.

The Makefile bundles all of this:

- `make db-up` runs `postgres:16` in Docker with `POSTGRES_USER=$(whoami)` and `POSTGRES_PASSWORD=postgres`, mapped to `localhost:$(DB_PORT)`. **Default `DB_PORT` is `54329`** (obscure on purpose to avoid clashing with other Postgres on 5432). Override with `DB_PORT=5432 make db-up`.
- `make test` exports `PGPORT=$(DB_PORT)` so tests find the container automatically.
- A committed `.envrc` sets `PGHOST=localhost`, `PGPORT=54329`, `PGPASSWORD=postgres` for direnv users — `psql` and ad-hoc `go test` outside Make pick up the same wiring. Run `direnv allow` once after pulling.

External tooling required: `sqlc` (v1.31.x), `mockery` (v3.7.x — the `.mockery.yaml` is v3 schema; v2 will not work), `direnv` (optional but recommended), and Docker. All available via `brew install`.

## Architecture: DDD-lite / threedots.tech layout

The application code lives under `pkg/withoutings/` and is split into four layers. **Imports flow inward only.** This is the most important rule when changing or adding code:

```
port  ─┐
       ├──►  app  ──►  domain
adapter┘            ▲
                    │
                  (adapter implements domain.Repo interfaces)
```

- **`domain/<aggregate>/`** — pure business types and behavior. Aggregates currently are `account`, `subscription`, `withings`. Defines entities, value objects, the `Repo` interface for that aggregate, and domain events (`subscription_events.go`). Domain packages MUST NOT import `app`, `adapter`, `port`, `db`, `web`, or any infrastructure package. They may import other domain packages and stdlib/uuid/etc.
- **`app/`** — application layer. Orchestrates domain via CQRS:
  - `app/command/` — write-side handlers (`SubscribeAccount`, `ProcessRawNotification`, `RefreshAccessToken`, `FetchNotificationData`, ...). Each handler is an interface + struct constructor; the struct depends on domain `Repo` interfaces and domain services, never on adapter types.
  - `app/query/` — read-side handlers (`AccountByUUID`, `AllAccounts`, ...).
  - `app/service/<name>/` — domain/application services that wrap external systems (e.g. `withings` service wraps the Withings API client and adds token-refresh logic). Inputs/outputs are domain types.
  - `app/app.go` is the composition root. It is the only place in `app/` that imports `adapter/*`. `App` struct holds repos, services, `Commands{}`, `Queries{}`, watermill `Publisher`/`Subscriber`. `App.Validate()` is called from `web.Router` to fail fast if a wire-up is missing.
- **`adapter/<aggregate>/`** — concrete implementations of the domain `Repo` interfaces.
  - Storage adapters (`account`, `subscription`) depend on `pkg/db` (sqlc-generated) and the matching `domain/<aggregate>` package.
  - The `withings` adapter is an HTTP client for the Withings API, not a Postgres adapter — it implements `domain/withings.Repo` and does not touch `pkg/db`.
  - `adapter/topic/` is a special non-aggregate package that just holds watermill topic-name constants; it is allowed to be imported by `port/` and `worker/` so they can publish/subscribe to the same string keys.
- **`port/`** — inbound HTTP handlers. Each file is a `func XxxPage(svc *app.App) http.HandlerFunc`. Handlers parse requests, call `svc.Commands.X.Handle(...)` or `svc.Queries.X.Handle(...)`, then render via `svc.Templates`. Ports may import `adapter/topic` for publishing event names but MUST NOT import storage adapters or `pkg/db` directly.

### Adding a new feature, by layer

1. Model it in `domain/<aggregate>/`. Add behaviour to the entity; extend `Repo` interface only if persistence semantics actually change.
2. Implement the `Repo` change in `adapter/<aggregate>/` and add the SQL in `pkg/db/<aggregate>.sql`. Run `make generate-sqlc`.
3. Add a `command/` or `query/` handler that depends on the domain `Repo` and any needed `app/service`. Wire it into `app.Commands` / `app.Queries` in `app/app.go` and add it to the `Validate()` checks.
4. Mount an HTTP handler in `port/` and route it in `pkg/web/routes.go`.
5. For async work, publish to a topic from `adapter/topic/topics.go` via `app.Publisher`; consume in `pkg/worker/worker.go`, which dispatches to a command handler.

### Async / event flow

- Pub/sub uses watermill with the **Postgres SQL** adapter (no Redis). `app.NewApplication` creates both a `*pgxpool.Pool` (used by everything) and a `*sql.DB` (required by `watermill-sql/v3`). Don't replace the SQL adapter without updating both.
- Topic names live in `adapter/topic/topics.go`. Domain events live next to the aggregate in `domain/<aggregate>/<aggregate>_events.go`. There is a known TODO to separate domain events from JSON/watermill encoding — keep new event payloads minimal (UUID references) until that split happens.
- `pkg/worker/worker.go` is the consumer. It subscribes per topic and routes messages into command handlers; it embeds `*app.App`, so adding a new consumer means adding a goroutine in `Work(ctx)`.

### Persistence

- Schema lives in `pkg/migration/*.sql` (golang-migrate, embedded). New migrations: append `NNNNNN_name.up.sql` (and `.down.sql` if reversible).
- Queries are sqlc-generated from `pkg/db/account.sql` and `pkg/db/subscription.sql`. Schema for sqlc is the migration directory. Regenerate via `make generate-sqlc` after schema or query changes.
- Sessions use `alexedwards/scs` with a Postgres store; do not introduce Redis.
- Repo interfaces support transactional composition via a `WithTx(pgx.Tx)` method on the adapter (see `accountAdapter.PgRepo.WithTx`). Multi-aggregate writes should run inside a single pgx transaction rather than across separate command handlers.

### Mocks and testing

- Mocks are mockery-generated in-package with `_mock` suffix and `EXPECT()` helpers (`.mockery.yaml` is authoritative). After adding/changing an interface marked with `//go:generate mockery ...`, run `make generate-all`.
- `pkg/integrationtest` brings up a fresh Postgres database and a fully wired `*app.App` (via `app.NewTestApplication`) plus a `*worker.Worker` and `mux.Router`. Use it for handler/worker tests that span layers.
- The binary entry point is `cmd/withoutings/{main,cli,migrate}.go`. There is no top-level `cmd/main.go`. The single binary dispatches subcommands `withoutings server` and `withoutings migrate`.

## Project conventions worth knowing

- No JS framework (`pkg/web/package.json` only pulls in tailwind/postcss/autoprefixer as dev tooling), no ORM, no first-party Redis use (a transitive `redis/go-redis` dep exists via watermill, but nothing in this repo depends on Redis at runtime). Server-side rendering via `pkg/web/templates`. Adding any of those should be discussed first.
- The `webhook secret` is part of the public callback URL path (`/withings/webhooks/<secret>`); treat the URL itself as a credential.
- Don't modify primary keys manually — sequences will conflict on next insert (see `pkg/db/README.md`).
