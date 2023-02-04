.PHONY: sqlc
sqlc:
	cd pkg/db && sqlc generate
