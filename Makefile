.PHONY: sqlc
sqlc:
	cd pkg/db && sqlc generate

.PHONY: run-dev
run-dev:
	source env.dev.sh && go run cmd/withoutings/*.go migrate && go run cmd/withoutings/*.go server