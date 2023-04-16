postgres:
	docker run --name pure_bank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15.2-alpine
createdb:
	docker exec -it pure_bank createdb --username=root --owner=root pure_bank

new_migration:
	migrate create -ext sql -dir db/migration/ -seq $(name)

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pure_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/pure_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

.PHONY:postgres createdb new_migration migrateup migratedown sqlc test
