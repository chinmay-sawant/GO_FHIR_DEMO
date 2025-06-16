# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=fhir-patient-api
BINARY_WINDOWS=$(BINARY_NAME).exe

# Database parameters
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: all build clean test deps run help migrate-up migrate-down migrate-create mocks

all: test build

## Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

## Build for Windows
build-windows:
	set GOOS=windows&& set GOARCH=amd64&& $(GOBUILD) -o $(BINARY_WINDOWS) -v ./...

## Clean build files
clean:
	$(GOCLEAN)
	del $(BINARY_NAME) 2>nul || echo.
	del $(BINARY_WINDOWS) 2>nul || echo.

## Run tests
test:
	$(GOTEST) -v ./...

## Download dependencies
deps:
	$(GOMOD) tidy
	$(GOMOD) download

## Run the application
run:
	$(GOCMD) run main.go

## Run database migrations up
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

## Run database migrations down
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

## Create a new migration file
migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

## Install migration tools
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## Setup development environment
setup: deps install-migrate
	copy .env.example .env
	echo Please update .env file with your database credentials


# Target to generate all mocks
.PHONY: mocks
mocks:
	@echo Starting mock generation...
	@for /f "tokens=*" %%f in ('dir /s /b "internal\*.go" 2^>nul ^| findstr /v /i "mock" ^| findstr /v /i "_test"') do ( \
		echo Processing Go file: %%f && \
		for %%d in ("%%~dpf") do ( \
			set "file_dir=%%~d" && \
			set "file_name=%%~nf" && \
			echo File directory: %%~d && \
			echo File name: %%~nf && \
			echo %%~nf | findstr /I "impl" >nul && ( \
				echo Skipping %%f because filename contains 'impl' \
			) || ( \
				if not exist "%%~dpfmocks" ( \
					echo Creating mocks directory: %%~dpfmocks && \
					mkdir "%%~dpfmocks" \
				) && \
				echo Command: mockgen -source=%%f -destination=%%~dpfmocks\mock_%%~nf.go -package=mocks && \
				mockgen -source="%%f" -destination="%%~dpfmocks\mock_%%~nf.go" -package=mocks \
			) \
		) \
	)
	@echo Mock generation completed.

# Clean target to remove all generated mocks
.PHONY: clean-mocks
clean-mocks:
	@echo Cleaning all mock files...
	@for /f "tokens=*" %%d in ('dir /s /b /ad "internal\*mocks" 2^>nul') do ( \
		echo Cleaning mock directory: %%d && \
		del /Q "%%d\*.go" 2>nul || echo No mocks to clean in %%d. \
	)
	@echo Mock cleanup completed.

## Display help
help:
	@echo Available commands:
	@echo   build         - Build the application
	@echo   build-windows - Build for Windows
	@echo   clean         - Clean build files
	@echo   test          - Run tests
	@echo   deps          - Download dependencies
	@echo   run           - Run the application
	@echo   migrate-up    - Run database migrations up
	@echo   migrate-down  - Run database migrations down
	@echo   migrate-create- Create new migration (use: make migrate-create name=migration_name)
	@echo   setup         - Setup development environment
	@echo   help          - Display this help
	@echo   mocks         - Generate all mocks
	@echo   clean-mocks   - Clean all generated mocks