.PHONY: all
all: generate-all test build

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: build
build:
	go build -v -o . ./...

.PHONY: generate-all
generate-all:
	go generate ./...

.PHONY: generate-sqlc
generate-sqlc:
	go generate ./pkg/db

.PHONY: migrate
migrate:
	source env.dev.sh && go run cmd/withoutings/*.go migrate

.PHONY: run-dev
run-dev: generate-all migrate
	source env.dev.sh && go run cmd/withoutings/*.go server