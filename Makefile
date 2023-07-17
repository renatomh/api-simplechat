DB_URL=postgresql://root:secret@localhost:5432/simple_chat?sslmode=disable

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_chat

dropdb:
	docker exec -it postgres12 dropdb simple_chat

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	docker run --rm -v "$(CURDIR):/src" -w //src kjconroy/sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
