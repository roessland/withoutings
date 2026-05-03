.PHONY: all
all: generate-all test build

.PHONY: test
test:
	PGPORT=$(DB_PORT) go test -v -race -cover ./...

.PHONY: build
build:
	go build -v -o . ./...

.PHONY: generate-all
generate-all: generate-mocks
	go generate ./...

.PHONY: generate-mocks
generate-mocks:
	mockery

.PHONY: generate-sqlc
generate-sqlc:
	go generate ./pkg/db

.PHONY: migrate
migrate:
	source env.dev.sh && go run cmd/withoutings/*.go migrate

.PHONY: run-dev
run-dev: generate-all migrate
	source env.dev.sh && go run cmd/withoutings/*.go server

.PHONY: clean
clean:
	go clean -testcache

# Local Postgres via Docker. Obscure port to avoid clashing with anything on 5432.
DB_CONTAINER := withoutings-db
DB_VOLUME    := withoutings-db-data
DB_IMAGE     := postgres:16
DB_PORT      ?= 54329

.PHONY: db-up
db-up:
	@docker start $(DB_CONTAINER) >/dev/null 2>&1 || \
		docker run -d --name $(DB_CONTAINER) \
			-e POSTGRES_USER=$$(whoami) \
			-e POSTGRES_PASSWORD=postgres \
			-p $(DB_PORT):5432 \
			-v $(DB_VOLUME):/var/lib/postgresql/data \
			$(DB_IMAGE) >/dev/null
	@until docker exec $(DB_CONTAINER) psql -U $$(whoami) -d postgres -c 'select 1' >/dev/null 2>&1; do sleep 0.2; done
	@echo "Postgres ready on localhost:$(DB_PORT) (superuser $$(whoami) / password postgres)"

.PHONY: db-down
db-down:
	@docker stop $(DB_CONTAINER) >/dev/null 2>&1 || true
	@docker rm $(DB_CONTAINER) >/dev/null 2>&1 || true
	@echo "Postgres container stopped (volume $(DB_VOLUME) retained)"

.PHONY: db-destroy
db-destroy: db-down
	@docker volume rm $(DB_VOLUME) >/dev/null 2>&1 || true
	@echo "Postgres volume $(DB_VOLUME) removed"

# Separate dev DB so test runs can be destroyed without losing dev data.
# Pre-bakes the wot database + wotsa/wotrw roles via deploy/dev/postgres-init.sql.
DEV_DB_CONTAINER := withoutings-dev-db
DEV_DB_VOLUME    := withoutings-dev-db-data
DEV_DB_PORT      ?= 54330

.PHONY: dev-db-up
dev-db-up:
	@docker start $(DEV_DB_CONTAINER) >/dev/null 2>&1 || \
		docker run -d --name $(DEV_DB_CONTAINER) \
			-e POSTGRES_USER=postgres \
			-e POSTGRES_PASSWORD=postgres \
			-p $(DEV_DB_PORT):5432 \
			-v $(DEV_DB_VOLUME):/var/lib/postgresql/data \
			-v $(CURDIR)/deploy/dev/postgres-init.sql:/docker-entrypoint-initdb.d/init.sql:ro \
			$(DB_IMAGE) >/dev/null
	@until docker exec $(DEV_DB_CONTAINER) psql -U wotsa -d wot -c 'select 1' >/dev/null 2>&1; do sleep 0.2; done
	@echo "Dev Postgres ready on localhost:$(DEV_DB_PORT) (db wot, roles wotsa/wotsa and wotrw/wotrw)"

.PHONY: dev-db-down
dev-db-down:
	@docker stop $(DEV_DB_CONTAINER) >/dev/null 2>&1 || true
	@docker rm $(DEV_DB_CONTAINER) >/dev/null 2>&1 || true
	@echo "Dev Postgres container stopped (volume $(DEV_DB_VOLUME) retained)"

.PHONY: dev-db-destroy
dev-db-destroy: dev-db-down
	@docker volume rm $(DEV_DB_VOLUME) >/dev/null 2>&1 || true
	@echo "Dev Postgres volume $(DEV_DB_VOLUME) removed"
