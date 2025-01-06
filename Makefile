migration:
	go run cmd/internal/db/migrate/migration.go

seed:
	go run cmd/internal/db/seeder/seeder.go

run:
	go run cmd/server/main.go

swagger:
	swag init -g cmd/server/main.go
