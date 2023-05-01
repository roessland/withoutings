# Withoutings

Demo application for talking with Withings API. It is written in Go with few dependenices, and uses PostgreSQL as a
database. The code is organized using DDD-ish / Clean Architecture-ish principles, based on the book and articles
by https://threedots.tech/.

## Features

### Currently available features

- Lets users log in with Withings OAuth and stores their access token in the database.
- Lets users subscribe to notifications from Withings.
- Stores received notifications in the database.

### Planned features

- Forward received notifications and their payloads to webhooks, IFTTT, Dropbox, Google Drive, etc.
- Download historical data, not just data corresponding to notifications.

# Installation and setup

The application serves a website and runs services that talk with the Withings API. Therefore it must have a public URL
that can receive webhooks sent by the Withings Notification service.

## Set up a Withings Developer account

Go to the [Withings Developer Dashboard](https://developer.withings.com/dashboard/).
Create a new application.

## Define environment variables

See [env.example.sh](env.example.sh). Make a copy of it named `env.dev.sh` and fill in the values.

```bash
source env.sh && go run cmd/main.go
```

The webhook secret must be added in the registered callback path in the Withings Developer Dashboard.
See `env.example.sh`.

# Development

## Forward remote port to local port

To receive webhooks in your development environment, you can forward a remote port to your local port.

Withings calls `https://withings.mywebsite.com/auth/callback` which is
forwarded to port 3628 on the server (e.g. using Caddy or nginx), which
is again forwarded to port 3628 in your development environment.

### Using SSH

```bash
# Using SSH
ssh -R 3628:127.0.0.1:3628 -N -f myuser@withings.mywebsite.com
```

### Using Caddy and Tailscale

Set up Tailscale on your development machine and the server. Then add the following to your Caddyfile:

```Caddyfile
withings-dev.example.com {
        reverse_proxy /* <dev-machine-name>:3628 {
        }
}
```
The server must also listen on the Tailscale interface. Configure that in `env.dev.sh`.
```shell
export WOT_LISTEN_ADDR='<dev-machine-tailscale-ip>:3628';
```

## Migrations

Migrations are managed using golang-migrate. The library is embedded in the build, so you can run migrations using `withoutings migrate`.

### Create migration

Append a new migration file in the `migration` directory.

### Run all necessary migrations

```
source env.sh && withoutings migrate
```

### Revert migration

Manually revert by executing the `down` SQL from the migration file.
Remember to also decrement the migration version in the `schema_migrations` table.

## SQL queries

Go code is generated from SQL queries using [sqlc](https://docs.sqlc.dev/).
The schema is inferred using the migration files.

### Install sqlc

```sh
brew install sqlc
# or
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

## Mock generation

Mocks are generated using [mockery](https://github.com/vektra/mockery).

### Install [mockery](https://github.com/vektra/mockery)

```sh
brew install mockery
```

### Generate mocks

For now, you have to generate mocks using go generate.