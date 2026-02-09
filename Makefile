include ./.env

export CGO_CFLAGS = -IC:\msys64\mingw64\include
export CGO_LDFLAGS = -LC:\msys64\mingw64\lib -ltesseract -lleptonica
export CC = C:\msys64\mingw64\bin\gcc.exe
export CXX = C:\msys64\mingw64\bin\g++.exe
export PATH := C:\msys64\mingw64\bin;$(PATH)

define db_up
	migrate -path ./migrations -database "sqlite3://${DATABASE_PATH}" up
endef

define db_down
	migrate -path ./migrations -database "sqlite3://${DATABASE_PATH}" down
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

prebuild:
	@echo Environment configured.

dev: prebuild
	@air

build:
	@go build -o bin/main ./cmd

run: build
	@bin/main