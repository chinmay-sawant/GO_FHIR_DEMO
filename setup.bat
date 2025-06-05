@echo off
echo Setting up Go FHIR Demo Application...

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH
    exit /b 1
)

REM Install dependencies
echo Installing Go dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo Error: Failed to install dependencies
    exit /b 1
)

REM Install migration tool
echo Installing migration tool...
go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
if %errorlevel% neq 0 (
    echo Error: Failed to install migration tool
    exit /b 1
)

REM Copy environment file
if not exist .env (
    echo Creating .env file from example...
    copy .env.example .env
    echo Please update .env file with your database credentials
) else (
    echo .env file already exists
)

echo.
echo Setup complete! 
echo.
echo Next steps:
echo 1. Update .env file with your PostgreSQL database credentials
echo 2. Ensure PostgreSQL is running
echo 3. Run database migrations: migrate -path migrations -database "your_db_url" up
echo 4. Start the application: go run main.go
echo.
echo For more commands, see the Makefile or run: make help
