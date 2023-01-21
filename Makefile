.PHONY: sqlc
sqlc:
	cd pkg/repos && sqlc generate
