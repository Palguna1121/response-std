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
	@echo "🔼 Running migration UP for version..."
	go run app/console/cmd/migrate/migrate.go up

migrate-down:
	@echo "🔽 Running migration DOWN for version..."
	go run app/console/cmd/migrate/migrate.go down

migrate-drop:
	@echo "💥 Dropping database for version..."
	go run app/console/cmd/migrate/migrate.go drop

migrate-force:
	@echo "⚙️ Forcing migration version on..."
	go run app/console/cmd/migrate/migrate.go force
#all migrate usage: make migrate-up/migrate-down/migrate-drop
#example: make migrate-up


# ================================================================================
# USING MIGRATE AND MODEL AT THE SAME TIME
# ================================================================================
# migrate and model
migrate-m:
	@echo "Running migration and generating models..."
# todo: run create migration and generate model


# ================================================================================
# ================================================================================
# ================================================================================


db-seed:
	go run app/console/cmd/seed/seed.go
#usage: make db-seed