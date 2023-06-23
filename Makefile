# USED ONLY IN DEVELOPMENT
export DB_URL=postgres://postgres:secret@localhost:5432/postgres?sslmode=disable
export PORT=3000
export DB_DRIVER=postgres
MIGRATION_NAME ?= new_migration

db-new-migration:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate create -ext sql -dir db/migrations -seq "$(MIGRATION_NAME)"

db-migrate-up:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate -path=./db/migrations -database "${DB_URL}" -verbose up 1
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate -path=./db/migrations -database "${DB_URL}" -verbose up 1

db-migrate-up-linux:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

db-migrate-down:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate -path=./db/migrations -database "${DB_URL}" -verbose down 1

generate-models: 
	docker run --rm -v "$(shell echo %cd%):/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

start:
	go run .
