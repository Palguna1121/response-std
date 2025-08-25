run:
	go run ./main.go

build:
	go build -o bin/app ./main.go

test:
	go test ./...

install:
	@echo "Installing dependencies..."
	go mod tidy

help:
	@echo "Available commands:"
	@echo "  run                               - Run the application"
	@echo "  build                             - Build the application"
	@echo "  test                              - Run tests"
	@echo "  install                           - Install dependencies"
	@echo "  help                              - Show this help message"
	@echo "  model name=NAME                   - Generate model for the specified NAME"
	@echo "  controller name=NAME ver=VERSION  - Generate controller for the specified NAME and VERSION"
	@echo "  request name=NAME ver=VERSION     - Generate request for the specified NAME and VERSION"
	@echo "  scaffold name=NAME ver=VERSION    - Generate model, controller, and requests for the specified NAME and VERSION"
	@echo "  scaffolds name=NAME ver=VERSION   - Create migration, generate model, controller, and requests for the specified NAME and VERSION"
	@echo "  create migration=MIGRATION_NAME   - Create a new migration file with the specified MIGRATION_NAME"
	@echo "  migrate-up                        - Run database migrations up"
	@echo "  migrate-down                      - Run database migrations down"
	@echo "  migrate-drop                      - Drop the database"
	@echo "  migrate-force                     - Force database migration to a specific version"
	@echo "  db-seed                           - Seed the database with initial data"
	@echo "  seeder name=NAME                  - Generate a seeder for the specified NAME"
	@echo "  fresh-seed                        - Fresh migrate and seed the database"
	@echo "  migrate-up-seed                   - Migrate up and seed the database"

# ================================================================================
# ================================================================================
# ================================================================================

# generate models
model:
	@echo "Generating models... $(name)"
	go run app/console/cmd/scripts/models/generate_models.go $(name)
#usage: make models name=users

#required migration: table for automation of model generation

#example:  
#	command: "make models name=users" 
#		it will generate models from migration database/migrations/create_users_table.up.sql
# 		thats why this command requires migration first

# ================================================================================
# ================================================================================
# ================================================================================


# generate controllers
controller:
	@echo "Generating controllers... $(name) (version: $(ver))"
	go run app/console/cmd/scripts/controllers/generate_controllers.go $(name) $(ver)
#usage: make controllers name=users ver=v1


# ================================================================================
# ================================================================================
# ================================================================================

# Generate requests
request:
	@echo "Generating requests... $(name) (version: $(ver))"
	go run app/console/cmd/scripts/requests/generate_requests.go $(name) $(ver)
# usage: make requests name=user_store ver=v1

# ================================================================================
# ================================================================================
# ================================================================================

# Generate all: models, controller and request
scaffold:
	@echo "Generating scaffold... $(name) (version: $(ver))"
	make model name=$(name)
	make controller name=$(name) ver=$(ver)
	make request name=$(name)_store ver=$(ver)
	make request name=$(name)_update ver=$(ver)
# usage: make scaffold name=users ver=v1

# ================================================================================
# ================================================================================
# ================================================================================

# Generate all: models, controller and request
scaffolds:
	@echo "Generating scaffold... $(name) (version: $(ver))"
	make create migration=create_$(name)_table
	make model name=$(name)
	make controller name=$(name) ver=$(ver)
	make request name=$(name)_store ver=$(ver)
	make request name=$(name)_update ver=$(ver)
# usage: make scaffolds name=users ver=v1

# ================================================================================
# ================================================================================
# ================================================================================

# migrate using golang-migrate/migrate driver mysql
MIGRATION_TOOL = migrate
MIGRATIONS_DIR = database/migrations

create:
	@echo "📝 Creating migration: $(migration) in $(MIGRATIONS_DIR)"
	@mkdir -p $(MIGRATIONS_DIR)
	$(MIGRATION_TOOL) create -ext sql -dir $(MIGRATIONS_DIR) $(migration)
#usage: make create migration=create_users_table

# ================================================================================
# ================================================================================
# ================================================================================


# command aliases for migrations
migrate-up:
	@echo "🔼 Running migration UP..."
	go run app/console/cmd/migrate/migrate.go up

migrate-down:
	@echo "🔽 Running migration DOWN..."
	go run app/console/cmd/migrate/migrate.go down

migrate-drop:
	@echo "💥 Dropping database..."
	go run app/console/cmd/migrate/migrate.go drop

migrate-force:
	@echo "⚙️ Forcing migration..."
	go run app/console/cmd/migrate/migrate.go force
#all migrate usage: make migrate-up/migrate-down/migrate-drop
#example: make migrate-up


# ================================================================================
# ================================================================================
# ================================================================================


db-seed:
	go run app/console/cmd/scripts/seed/run.go
#usage: make db-seed

seeder:
	@echo "Generating Seeder... $(name)"
	go run app/console/cmd/scripts/seed/generate/generate_seeder.go $(name)


fresh-seed:
	@echo "Running fresh migration and seeding database..."
	make migrate-down
	make migrate-up
	make db-seed

migrate-up-seed:
	@echo "Running migration UP and seeding database..."
	make migrate-up
	make db-seed