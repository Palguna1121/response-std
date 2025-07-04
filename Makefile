run:
	go run ./main.go

build:
	go build -o bin/app ./main.go

test:
	go test ./...

install:
	@echo "Installing dependencies..."
	go mod tidy

# ================================================================================
# ================================================================================
# ================================================================================

# generate models
model:
	@echo "Generating models... $(name) (version: $(ver))"
	go run cmd/scripts/models/generate_models.go $(name) $(ver)
#usage: make models name=users ver=v1

#required migration: table for automation of model generation

#example:  
#	command: "make models name=users ver=v1" 
#		it will generate models from migration v1/migrations/create_users_table.up.sql
# 		thats why this command requires migration first

# ================================================================================
# ================================================================================
# ================================================================================


# generate controllers
controller:
	@echo "Generating controllers... $(name) (version: $(ver))"
	go run cmd/scripts/controllers/generate_controllers.go $(name) $(ver)
#usage: make controllers name=users ver=v1


# ================================================================================
# ================================================================================
# ================================================================================

# Generate requests
request:
	@echo "Generating requests... $(name) (version: $(ver))"
	go run cmd/scripts/requests/generate_requests.go $(name) $(ver)
# usage: make requests name=user_store ver=v1

# ================================================================================
# ================================================================================
# ================================================================================

# Generate all: models, controller and request
scaffold:
	@echo "Generating scaffold... $(name) (version: $(ver))"
	make model name=$(name) ver=$(ver)
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
	make create migration=create_$(name)_table ver=$(ver)
	make model name=$(name) ver=$(ver)
	make controller name=$(name) ver=$(ver)
	make request name=$(name)_store ver=$(ver)
	make request name=$(name)_update ver=$(ver)
# usage: make scaffolds name=users ver=v1

# ================================================================================
# ================================================================================
# ================================================================================

# migrate using golang-migrate/migrate driver mysql
# default version fallback = v1
ver ?= v1
MIGRATION_TOOL = migrate
MIGRATIONS_DIR = $(ver)/database/migrations

create:
	@echo "📝 Creating migration: $(migration) in $(MIGRATIONS_DIR)"
	@mkdir -p $(MIGRATIONS_DIR)
	$(MIGRATION_TOOL) create -ext sql -dir $(MIGRATIONS_DIR) $(migration)
#usage: make create migration=create_users_table ver=v2

# ================================================================================
# ================================================================================
# ================================================================================


# command aliases for migrations
migrate-up:
	@echo "🔼 Running migration UP for version $(ver)..."
	go run cmd/migrate/migrate.go up $(ver)

migrate-down:
	@echo "🔽 Running migration DOWN for version $(ver)..."
	go run cmd/migrate/migrate.go down $(ver)

migrate-drop:
	@echo "💥 Dropping database for version $(ver)..."
	go run cmd/migrate/migrate.go drop $(ver)

migrate-force:
	@echo "⚙️ Forcing migration version $(VERSION) on $(ver)..."
	go run cmd/migrate/migrate.go force $(ver) $(VERSION)
#all migrate usage: make migrate-up/migrate-down/migrate-drop then version(ver=v1)
#example: make migrate-up ver=v2

# ================================================================================
# ================================================================================
# ================================================================================


db-seed:
	go run cmd/seed/seed.go
#usage: make db-seed