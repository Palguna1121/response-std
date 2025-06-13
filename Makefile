run:
	go run ./main.go

build:
	go build -o bin/app ./main.go

test:
	go test ./...


# migrate using golang-migrate/migrate driver mysql
migrate-up:
	@echo "Running migration UP..."
	go run cmd/migrate/migrate.go up

migrate-down:
	@echo "Running migration DOWN..."
	go run cmd/migrate/migrate.go down

migrate-drop:
	@echo "Dropping database..."
	go run cmd/migrate/migrate.go drop

migrate-force:
	@echo "Forcing version $(VERSION)..."
	go run cmd/migrate/migrate.go force $(VERSION)

db-seed:
	go run cmd/seed/seed.go
