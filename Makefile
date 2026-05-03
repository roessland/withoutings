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
