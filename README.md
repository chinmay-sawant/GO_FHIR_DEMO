# Go FHIR Demo Application

A Go Gin framework application with FHIR (Fast Healthcare Interoperability Resources) support using PostgreSQL database.

## Features

- RESTful API for Patient resources (GET, POST, PUT, DELETE, PATCH)
- FHIR R4 compliance using official FHIR Go library
- PostgreSQL database with GORM
- Database migrations using [golang-migrate](https://github.com/golang-migrate/migrate) (see below for Windows install instructions)
- Structured logging with Logrus
- Configuration management with Viper
- Request/Response time tracking middleware
- Clean architecture with interfaces

## Project Structure

```
├── config/           # Configuration files
├── internal/
│   ├── api/         # HTTP handlers and routes
│   ├── domain/      # Domain models and interfaces
│   ├── repository/  # Data access layer
│   ├── service/     # Business logic layer
│   └── middleware/  # Custom middleware
├── pkg/             # Shared packages
├── migrations/      # Database migration files
├── logs/           # Log files
└── main.go         # Application entry point
```

## Setup

## Swagger/OpenAPI Documentation

After running the application, you can access the interactive Swagger UI at:

    http://localhost:8080/swagger/index.html

This provides a browsable interface for all API endpoints and models.


1. Install dependencies:
```bash
go mod tidy
```

2. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_username
export DB_PASSWORD=your_password
export DB_NAME=fhir_demo
export GIN_MODE=debug
export LOG_LEVEL=info
```


3. Install golang-migrate (migration tool):

**On Windows:**
1. Download the latest Windows release from: https://github.com/golang-migrate/migrate/releases
2. Extract `migrate.exe` and place it in a folder (e.g., `C:\tools\migrate.exe`).
3. Add that folder to your system PATH (so you can run `migrate` from any terminal).
4. Open a new terminal and run:
   ```cmd
   migrate --version
   ```
   You should see the version output.

4. Run database migrations:
```cmd
migrate -path migrations -database "postgres://fhir_user:fhir_password@localhost:5432/fhir_demo?sslmode=disable" up
migrate -path migrations -database "postgres://username:password@localhost:5432/fhir_demo?sslmode=disable" up
```

4. Run the application:
```bash
go run main.go
```

5. (Optional) Regenerate Swagger docs after editing annotations:
```bash
swag init
```

## API Endpoints

- `GET /api/v1/patients` - Get all patients
- `GET /api/v1/patients/:id` - Get patient by ID
- `POST /api/v1/patients` - Create new patient
- `PUT /api/v1/patients/:id` - Update patient
- `PATCH /api/v1/patients/:id` - Partially update patient
- `DELETE /api/v1/patients/:id` - Delete patient

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | - |
| DB_PASSWORD | Database password | - |
| DB_NAME | Database name | - |
| GIN_MODE | Gin mode (debug/release) | debug |
| LOG_LEVEL | Log level | info |
| SERVER_PORT | Server port | 8080 |
