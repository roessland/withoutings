# Withoutings

Demo application for talking with Withings API.

## Setup

Go to the [Withings Developer Dashboard](https://developer.withings.com/dashboard/).
Create a new application.

## Required environment variables

See [env.example.sh](env.example.sh).

```bash
source env.sh && go run cmd/main.go
```

The webhook secret must be added in the registered callback path in the Withings Developer Dashboard. See `env.example.sh`.

## Forward remote port to local port

To make callback to a remote server call your development app you can
use SSH port forwarding.

```bash
ssh -R 3628:127.0.0.1:3628 -N -f myuser@withings.mywebsite.com
```

So Withings calls `https://withings.mywebsite.com/auth/callback` which is
forwarded to port 3628 on the server (e.g. using Caddy or nginx), which
is again forwarded to port 3628 in your development environment.


## Migrations


### Install golang-migrate locally and remotely

```
#Linux:
curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz

#MacOS:
brew install golang-migrate
```

### Create migration
```
migrate create -ext sql -dir migrations -seq <migration_name>
```

### Run migrations
```
source env.sh && withoutings migrate
```

### Revert migration
```sh
migrate down
```

## Queries
Go code is generated from SQL queries using [sqlc](https://docs.sqlc.dev/)

### Install sqlc
```sh
brew install sqlc
# or
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

## Generate mocks
### Install [mockery](https://github.com/vektra/mockery)
```sh
brew install mockery
```

### Generate mocks
```sh
make mockery
```