.PHONY: sqlc
sqlc:
	cd pkg/db && sqlc generate

.PHONY: mockery
mockery:
	cd pkg/withoutings/domain/withings && mockery --all --inpackage --case underscore --with-expecter