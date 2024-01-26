postgres:
	docker run --name my-postgres --network bank-network -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres
mysql:
	docker run --name my-mysql --network bank-network -p 3306:3306 -e MYSQL_ROOT_PASSWORD=12345678 -d mysql
createdb:
	docker exec -it my-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it my-postgres dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down


migrateup1:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server