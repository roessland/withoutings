.PHONY: sqlc
sqlc:
	cd internal/repos && sqlc generate
