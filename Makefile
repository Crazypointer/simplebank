postgres:
	docker run --name my-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
createdb:
	docker exec -it my-postgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it my-postgres dropdb --username=postgres simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc