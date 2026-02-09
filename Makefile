include ./.env

define db_up
	    migrate -path ./migrations -database "sqlite3://${DATABASE_PATH}"  up
endef

define db_down
	    migrate -path ./migrations -database "sqlite3://${DATABASE_PATH}"  down
endef

define db_force
	    migrate -database "sqlite3://${DATABASE_PATH}" -path ./migrations force $(version)
endef

define db_create
	migrate create -ext sql -dir ./migrations -seq $(migration)
endef

db_force:
	$(call db_force)

db_down:
	$(call db_down)

db_up:
	$(call db_up)

db_create:
	$(call db_create)

dev:
	@refresh

build:
	@go build -o bin/main ./cmd

run: build
	@bin/main