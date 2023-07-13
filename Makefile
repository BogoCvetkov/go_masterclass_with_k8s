# USED ONLY IN DEVELOPMENT
MIGRATION_NAME ?= new_migration

db-new-migration:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate create -ext sql -dir db/migrations -seq "$(MIGRATION_NAME)"

db-migrate-up:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate -path=./db/migrations -database "${DB_URL}" -verbose up 

db-migrate-up-linux:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up
	echo "DB migrated"

db-migrate-down:
	docker run --rm -v "$(shell echo %cd%):/src" -w /src --network host migrate/migrate -path=./db/migrations -database "${DB_URL}" -verbose down 1

generate-models: 
	docker run --rm -v "$(shell echo %cd%):/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

start-api:
	go run cmd\api\main.go

start-workers:
	go run cmd\async\main.go

proto-generate:
	del /Q /F .\pb\*.go >nul 2>&1
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out ./pb  --grpc-gateway_opt logtostderr=true  --grpc-gateway_opt paths=source_relative  --grpc-gateway_opt generate_unbound_methods=true \
	--openapiv2_out=use_go_templates=true:swagger  --openapiv2_opt allow_merge=true,merge_file_name=masterclass   ./proto/*.proto
	statik -src=swagger -f
