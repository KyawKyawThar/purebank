DB_URL=postgresql://root:secret@localhost:5432/pure_bank?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.2-alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root pure_bank

new_migration:
	migrate create -ext sql -dir db/migration/ -seq $(name)

migrateup:
	migrate -path db/migration -database $(DB_URL) -verbose up

migratedown:
	migrate -path db/migration -database $(DB_URL) -verbose down

migrateup1:
	migrate -path db/migration -database $(DB_URL) -verbose up 1

migratedown1:
	migrate -path db/migration -database $(DB_URL) -verbose down 1

sqlc:
	sqlc generate

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

test:
	go test -v -cover -short ./...

#purebank/db get from go.mod module sometime it can be github.com/KKT/purebank
mock:
	mockgen -package mockdb -destination db/mock/store.go purebank/db/sqlc Store

.PHONY:postgres createdb new_migration migrateup migratedown migrateup1 migratedown1 sqlc test redis
