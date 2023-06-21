export POSTGRESQL_URL=postgres://postgres:secret@localhost:5432/postgres?sslmode=disable
MIGRATION_NAME ?= new_migration

db-new-migration:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate create -ext sql -dir db/migrations -seq "$(MIGRATION_NAME)"

db-migrate-up:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate -path=./db/migrations -database "${POSTGRESQL_URL}" -verbose up 1

db-migrate-up-linux:
	migrate -path db/migrations -database "$(POSTGRESQL_URL)" -verbose up

db-migrate-down:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate -path=./db/migrations -database "${POSTGRESQL_URL}" -verbose down 1

generate-models: 
	docker run --rm -v "$(shell echo %cd%):/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...
