# Go FHIR Demo Application

A comprehensive Go Gin framework application with FHIR (Fast Healthcare Interoperability Resources) R4 support, featuring PostgreSQL database, automatic API documentation with Swagger, external FHIR server integration, service discovery with Consul, secret management with Vault, and production-ready architecture.

## 🚀 Features

### Core Features
- **RESTful API** for Patient resources with full CRUD operations (GET, POST, PUT, PATCH, DELETE)
- **FHIR R4 Compliance** with standard FHIR data structures and validation
- **External FHIR Server Integration** - Connect to and query external FHIR servers (like HAPI FHIR)
- **FHIR Client Package** - Reusable HTTP client for external FHIR server communication
- **PostgreSQL Database** with GORM ORM and JSONB support for efficient FHIR data storage
- **Database Migrations** using [golang-migrate](https://github.com/golang-migrate/migrate)
- **Swagger/OpenAPI Documentation** with interactive UI and auto-generation
- **Automatic Data Seeding** with sample FHIR patient records on startup
- **Structured Logging** with configurable levels and formats
- **Configuration Management** with Viper supporting JSON files and environment variables
- **Request/Response Middleware** for performance monitoring, CORS, and error handling
- **Clean Architecture** with proper separation of concerns (handlers, services, repositories)

### Cron Job Handler

This project includes a custom **cron job handler** for simulating background jobs (such as data sync and cleanup) via API endpoints.

- **Endpoints:**
  - `POST /api/v1/cron/sync` - Triggers 5 background jobs with random delays (each job runs in a goroutine).
  - `POST /api/v1/cron/cleanup` - Cleans up jobs that are still in the "queued" state.

- **Implementation Details:**
  - **No external cron plugin is used.** The handler is implemented using Go's built-in `time`, `sync`, and goroutines for concurrency and job management.
  - Each job is tracked in a map with its state: `"queued"`, `"started"`, or `"completed"`.
  - A mutex (`sync.Mutex`) is used to ensure thread-safe access to the job map.

#### Why is it not a good idea to close existing running jobs?

Closing (cancelling) already running jobs in Go is not straightforward, especially if they are running in goroutines without explicit cancellation support (such as context cancellation). Forcefully terminating goroutines is not supported by the Go runtime and can lead to resource leaks or inconsistent state.

**Best Practice:**  
- Design jobs to be cancellable using `context.Context` so they can gracefully exit when requested.
- Only jobs that are still in the "queued" state (not yet started) can be safely cancelled/removed.

#### Cleanup Behavior in This Project

- **Job with ID `99`** is always in the `"queued"` state and can be closed/removed by the cleanup endpoint.
- **Jobs with IDs `1-5`** are started immediately in goroutines and move to the `"started"` state; these cannot be stopped once running and will complete naturally.
- When `/cron/cleanup` is called, only jobs in the `"queued"` state (like `99`) are removed. Running jobs (in goroutines) are not stopped.

#### Mutexes in Cron Job Handler
- This section uses a mutex within the cron job to ensure safe concurrent access to the shared object (h.jobs).
- The mutex is locked twice to prevent race conditions, as multiple goroutines may attempt to access or modify
- h.jobs at the same time. This approach guarantees thread safety when working with goroutines in the cron context.
---
### Service Discovery & Secret Management
- **Consul Integration** - Service discovery, service registration, and key-value store
- **HashiCorp Vault Integration** - Secure secret management and storage
- **Redis Cache** - High-performance caching for external FHIR API calls
- **Consul Handler** - API endpoint to fetch secrets from Consul KV store
- **Service Registration** - Automatic registration with Consul on startup

### Technical Features
- Docker support with docker-compose
- Redis caching for improved performance
- Timeout handling for external API calls
- Environment-based configuration with `.env` file support
- JSONB storage for efficient FHIR data querying
- HTTP client with timeout and error handling for external FHIR servers
- Comprehensive error handling and validation
- Production-ready logging and monitoring

## 📊 Project Structure

```
├── config/                   # Configuration management
│   ├── config.go            # Configuration loader with Viper
│   └── config.json          # Default configuration
├── docs/                    # Auto-generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── examples/                # Sample FHIR data
│   └── sample_patient.json  # Example patient resource
├── internal/                # Private application code
│   ├── api/
│   │   ├── handlers/        # HTTP request handlers
│   │   │   ├── patient_handler.go              # Local patient CRUD operations
│   │   │   ├── external_patient_handler.go     # External FHIR server integration
│   │   │   ├── consul_handler.go               # Consul KV secret management
│   │   │   └── cron/                           # Cron job handlers
│   │   │        └── cronjob_handler.go         # Cron job API logic (sync/cleanup)
│   │   └── routes/          # Route definitions and middleware setup
│   │       └── routes.go
│   ├── domain/              # Domain models and business entities
│   │   ├── patient.go       # FHIR Patient domain model
│   │   └── external_patient.go  # External patient service interface
│   ├── middleware/          # HTTP middleware
│   │   └── middleware.go    # CORS, logging, timing, error handling
│   ├── repository/          # Data access layer
│   │   └── patient_repository.go  # PostgreSQL data operations
│   └── service/             # Business logic layer
│       ├── patient_service.go           # Local patient business logic
│       └── external_patient_service.go  # External FHIR server service
├── logs/                    # Application logs
├── migrations/              # Database schema migrations
│   ├── 000001_create_patients_table.up.sql
│   └── 000001_create_patients_table.down.sql
├── pkg/                     # Shared/reusable packages
│   ├── database/            # Database connection utilities
│   ├── fhirclient/          # HTTP client for external FHIR servers
│   ├── logger/              # Structured logging utilities
│   └── utils/               # Common utility functions
│       ├── consul.go        # Consul KV utilities
│       └── consul/          # Consul service registration
│           └── register.go
├── vault/                   # Vault configuration
│   └── config/
│       └── vault.hcl        # Vault server configuration
├── docker-compose.yml       # Multi-service Docker setup
├── Dockerfile              # Application container definition
├── Makefile                # Development automation scripts
├── INSTALLATION.md         # Detailed installation guide
└── main.go                 # Application entry point with service registration
```

## 🛠️ Technologies Used

### Backend Framework & Language
- **Go 1.23** - Programming language
- **Gin Web Framework** - HTTP web framework
- **GORM** - Object-relational mapping library

### Database & Storage
- **PostgreSQL** - Primary database with JSONB support
- **golang-migrate** - Database migration tool

### FHIR Integration
- **golang-fhir-models** - FHIR R4 data models
- **Custom FHIR Client** - HTTP client for external FHIR servers
- **FHIR R4 Compliance** - Full support for Patient resources

### Documentation & API
- **Swagger/OpenAPI 3.0** - API documentation
- **gin-swagger** - Swagger middleware integration
- **swaggo/swag** - Swagger documentation generator

### Configuration & Environment
- **Viper** - Configuration management
- **godotenv** - Environment variable loading
- **Logrus** - Structured logging

### Development & Deployment
- **Docker & Docker Compose** - Containerization
- **Makefile** - Build automation
- **Air** (optional) - Live reloading for development

### Service Discovery & Secret Management
- **Consul** - Service discovery and key-value store integration
- **HashiCorp Vault** - Secure secret management and storage

## 📋 Prerequisites

For detailed installation instructions, see [INSTALLATION.md](INSTALLATION.md).

- **Go 1.23+** - Programming language runtime
- **PostgreSQL 15+** - Database server
- **Docker & Docker Compose** - Containerization platform
- **Git** - Version control system

## 🚀 Quick Start

For complete installation and setup instructions, please refer to [INSTALLATION.md](INSTALLATION.md).

### Quick Docker Setup
```bash
# Clone the repository
git clone <repository-url>
cd Go_FHIR_Demo

# Start all services with Docker Compose
docker-compose up -d

# Access the application
# API: http://localhost:8080
# Swagger UI: http://localhost:8080/swagger/index.html
# Consul UI: http://localhost:8500
# Vault UI: http://localhost:8200 (requires initialization)
```

## 📖 API Documentation

### Swagger UI
After running the application, access the interactive Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

This provides a complete browsable interface for all API endpoints, request/response schemas, and allows you to test the API directly from the browser.

### Regenerating Documentation
After making changes to API annotations, regenerate the Swagger docs:
```cmd
swag init --parseDependency --parseDepth 99
```

## 🔗 API Endpoints

### Core Application Endpoints

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| `GET` | `/health` | Application health check | Service status and version |
| `GET` | `/metadata` | FHIR capability statement | Server capabilities and supported operations |
| `GET` | `/swagger/index.html` | Interactive API documentation | Swagger UI interface |
| `POST` | `/api/v1/cron/sync` | Trigger 5 background data sync jobs (each with random delay) | Job trigger status |
| `POST` | `/api/v1/cron/cleanup` | Cleanup jobs in "queued" state (e.g., job 99) | Cleanup status and cleaned job IDs |

### Local Patient Resource Endpoints (FHIR R4 Compliant)

| Method | Endpoint | Description | Request Body | Query Parameters |
|--------|----------|-------------|--------------|------------------|
| `GET` | `/api/v1/patients` | Get all patients with pagination | - | `limit` (default: 10), `offset` (default: 0) |
| `GET` | `/api/v1/patients/{id}` | Get patient by ID | - | - |
| `POST` | `/api/v1/patients` | Create new patient | FHIR Patient JSON | - |
| `PUT` | `/api/v1/patients/{id}` | Update entire patient resource | FHIR Patient JSON | - |
| `PATCH` | `/api/v1/patients/{id}` | Partially update patient | Partial updates map | - |
| `DELETE` | `/api/v1/patients/{id}` | Delete patient (soft delete) | - | - |

### External FHIR Server Endpoints

| Method | Endpoint | Description | Request Body | Query Parameters |
|--------|----------|-------------|--------------|------------------|
| `GET` | `/api/v1/external-patients/{id}` | Get patient from external FHIR server by ID | - | - |
| `GET` | `/api/v1/external-patients/{id}/cached` | Get patient with Redis caching | - | - |
| `GET` | `/api/v1/external-patients/{id}/delayed` | Get patient with timeout control | - | `timeout` (seconds, default: 10) |
| `GET` | `/api/v1/external-patients` | Search patients on external FHIR server | - | FHIR search params (`name`, `family`, `given`, `birthdate`, `gender`, etc.) |
| `POST` | `/api/v1/external-patients` | Create patient on external FHIR server | FHIR-compliant Patient JSON | - |

### Service Discovery & Secret Management Endpoints

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| `GET` | `/consul/secret` | Get secret from Consul KV store | JSON secret data |
| `GET` | `/vault/secret` | Get secret from Vault KV store | JSON secret data |

### Example Usage

#### Create a New Patient (Local)
```bash
curl -X POST http://localhost:8080/api/v1/patients \
  -H "Content-Type: application/json" \
  -d @examples/sample_patient.json
```

#### Get Patient from External FHIR Server
```bash
curl -X GET "http://localhost:8080/api/v1/external-patients/123"
```

#### Get Patient with Redis Cache
```bash
curl -X GET "http://localhost:8080/api/v1/external-patients/123/cached"
```

#### Get Patient with Custom Timeout
```bash
curl -X GET "http://localhost:8080/api/v1/external-patients/123/delayed?timeout=30"
```

#### Search External FHIR Server
```bash
curl -X GET "http://localhost:8080/api/v1/external-patients?family=Smith&gender=female"
```

#### Get All Local Patients with Pagination
```bash
curl -X GET "http://localhost:8080/api/v1/patients?limit=20&offset=0"
```

#### Get Patient by ID
```bash
curl -X GET http://localhost:8080/api/v1/patients/1
```

#### Update Patient
```bash
curl -X PUT http://localhost:8080/api/v1/patients/1 \
  -H "Content-Type: application/json" \
  -d @examples/sample_patient.json
```

#### Partially Update Patient
```bash
curl -X PATCH http://localhost:8080/api/v1/patients/1 \
  -H "Content-Type: application/json" \
  -d '{"family": "UpdatedLastName"}'
```

#### Delete Patient
```bash
curl -X DELETE http://localhost:8080/api/v1/patients/1
```

#### Get Secret from Consul
```bash
curl -X GET http://localhost:8080/consul/secret
```

#### Get Secret from Vault
```bash
curl -X GET http://localhost:8080/vault/secret
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Database host address | `localhost` | Yes |
| `DB_PORT` | Database port | `5432` | Yes |
| `DB_USER` | Database username | - | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | - | Yes |
| `DB_SSLMODE` | SSL mode for database connection | `disable` | No |
| `SERVER_PORT` | HTTP server port | `8080` | No |
| `GIN_MODE` | Gin framework mode (`debug`/`release`) | `debug` | No |
| `LOG_LEVEL` | Logging level (`trace`/`debug`/`info`/`warn`/`error`) | `info` | No |
| `EXTERNAL_FHIR_SERVER_BASE_URL` | Base URL for external FHIR server | - | Yes |
| `CONSUL_ADDRESS` | Consul server address | `http://localhost:8500` | No |
| `CONSUL_KEY` | Consul KV key to fetch | `myapp/secret` | No |
| `VAULT_ADDRESS` | Vault server address | `http://localhost:8200` | No |
| `VAULT_TOKEN` | Vault authentication token | `root` | No |
| `VAULT_SECRET_PATH` | Vault secret path | `secret/data/myapp` | No |
| `REDIS_HOST` | Redis server host | `localhost` | No |
| `REDIS_PORT` | Redis server port | `6379` | No |
| `REDIS_PASSWORD` | Redis password | `` | No |
| `REDIS_DB` | Redis database number | `0` | No |

### Configuration File
The application also supports JSON configuration via `config/config.json` for default values. Environment variables take precedence over configuration file settings.

#### External FHIR Server Configuration
Configure the base URL for the external FHIR server in your environment:
```env
EXTERNAL_FHIR_SERVER_BASE_URL=http://hapi.fhir.org/baseR4
```

#### Consul Configuration
Configure Consul integration in your environment:
```env
CONSUL_ADDRESS=http://localhost:8500
CONSUL_KEY=myapp/secret
```
Or in `config/config.json`:
```json
"consul": {
  "address": "http://localhost:8500",
  "key": "myapp/secret"
}
```

#### Vault Configuration
Configure Vault integration in your environment:
```env
VAULT_ADDRESS=http://localhost:8200
VAULT_TOKEN=root
VAULT_SECRET_PATH=secret/data/myapp
```
Or in `config/config.json`:
```json
"vault": {
  "address": "http://localhost:8200",
  "token": "root",
  "secret_path": "secret/data/myapp"
}
```

#### Redis Configuration
Configure Redis for caching in your environment:
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

Popular public FHIR servers for testing:
- **HAPI FHIR R4:** `http://hapi.fhir.org/baseR4`
- **SMART Health IT:** `https://r4.smarthealthit.org`

## 🗄️ Database Schema

### Patients Table
```sql
CREATE TABLE patients (
    id SERIAL PRIMARY KEY,
    fhir_data JSONB NOT NULL,           -- Complete FHIR Patient resource
    active BOOLEAN,                      -- Patient active status
    family VARCHAR(255),                 -- Family name (for indexing)
    given VARCHAR(255),                  -- Given name (for indexing)
    gender VARCHAR(20),                  -- Patient gender
    birth_date DATE,                     -- Date of birth
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete support
);
```

### Indexes
- Performance indexes on `active`, `family`, `given`, `gender`, `birth_date`
- GIN index on `fhir_data` JSONB column for efficient JSON querying
- Soft delete index on `deleted_at`

## 🧪 Database Migrations

### Available Migration Commands

```cmd
# Run all pending migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Rollback the last migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down 1

# Check migration status
migrate -path migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" version

# Create a new migration
migrate create -ext sql -dir migrations -seq add_new_table
```

### Migration Files
- `000001_create_patients_table.up.sql` - Creates the patients table with indexes
- `000001_create_patients_table.down.sql` - Drops the patients table

## 🔨 Makefile Usage

The project includes a comprehensive Makefile for common development tasks:

### Build Commands
```cmd
# Build the application
make build

# Build for Windows specifically
make build-windows

# Clean build artifacts
make clean
```

### Development Commands
```cmd
# Download dependencies
make deps

# Run the application
make run

# Run tests
make test

# Setup development environment (copies .env.example to .env)
make setup
```

### Mock Generation Commands
```cmd
# Generate all mocks
make mocks

# Clean all generated mocks
make clean-mocks
```

### Database Commands
```cmd
# Run database migrations up
make migrate-up

# Run database migrations down
make migrate-down

# Create a new migration file
make migrate-create name=your_migration_name

# Install migration tools
make install-migrate
```

### Help
```cmd
# Display all available commands
make help
```

**Note:** If you don't have `make` installed on Windows, you can install it via:
- **Chocolatey:** `choco install make`
- **Manual:** Download from [GnuWin32](http://gnuwin32.sourceforge.net/packages/make.htm)

## 🐳 Docker Deployment

### Using Docker Compose (Recommended)
```cmd
# Start all services (PostgreSQL + Consul + Application)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Manual Docker Build
```cmd
# Build the Docker image
docker build -t go-fhir-demo .

# Run the container
docker run -p 8080:8080 --env-file .env go-fhir-demo
```

### Docker Environment
The `docker-compose.yml` includes:
- PostgreSQL database with persistent volume
- Consul service for key-value store and service discovery
- Application service with environment variables
- Network configuration for service communication

## 🛠️ Development Guide

### Project Development Setup
1. **Fork and clone** the repository
2. **Install dependencies** with `go mod tidy`
3. **Setup environment** with `make setup`
4. **Configure database** connection in `.env`
5. **Run migrations** with `make migrate-up`
6. **Start development** with `make run`

### Code Structure Guidelines
- **Handlers** - HTTP request/response logic only
- **Services** - Business logic and orchestration
- **Repositories** - Data access layer
- **Domain** - Business entities and interfaces
- **Config** - Configuration management
- **Middleware** - Cross-cutting concerns

### Adding New Features
1. **Domain Model** - Define entities in `internal/domain/`
2. **Repository** - Add data access in `internal/repository/`
3. **Service** - Implement business logic in `internal/service/`
4. **Handler** - Add HTTP endpoints in `internal/api/handlers/`
5. **Routes** - Register routes in `internal/api/routes/`
6. **Migration** - Create database changes in `migrations/`
7. **Consul Integration** - Add handler in `internal/api/handlers/consul_handler.go`, utility in `pkg/utils/consul.go`, and register route in `internal/api/routes/`.
8. **Documentation** - Update Swagger annotations and regenerate docs

### Testing
```cmd
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Testing & Coverage Commands
```cmd
# Generate JUnit XML report with verbose output
make test-with-junit

# Generate test coverage report
make coverage

# Generate detailed test coverage report
make coverage-detailed

# Clean coverage files
make clean-coverage

# Generate all mocks for testing
make mocks

# Clean generated mocks
make clean-mocks
```

### Test Coverage

The application includes comprehensive test coverage generation through the Makefile:

#### JUnit XML Report Generation
```cmd
make test-with-junit
```
This command:
- Uses gotestsum for enhanced test output formatting
- Generates a JUnit XML report at `junit-report.xml`
- Provides standard verbose output for detailed test information
- Compatible with CI/CD systems for test result publishing

#### Basic Coverage Report
```cmd
make coverage
```
This command:
- Runs all tests with coverage profiling
- Generates an HTML coverage report at `coverage.html`
- Provides a quick overview of test coverage

#### Detailed Coverage Report
```cmd
make coverage-detailed
```
This command:
- Runs tests with verbose output and atomic coverage mode
- Displays detailed coverage statistics in the console
- Generates both console summary and HTML report
- Shows function-level coverage information

#### Coverage File Management
```cmd
make clean-coverage
```
This command removes generated coverage files (`coverage.out`, `coverage.html`, and `junit-report.xml`).

#### Test Report Locations
- **JUnit XML**: `junit-report.xml` - Test results in JUnit format for CI/CD systems
- **Console Output**: Function-level coverage summary (with `make coverage-detailed`)
- **HTML Report**: `coverage.html` - Interactive coverage visualization
- **Raw Data**: `coverage.out` - Coverage profile data

### CI/CD Integration
The JUnit XML reports can be easily integrated into various CI/CD pipelines:

#### GitHub Actions Example
```yaml
- name: Run Tests
  run: make test-with-junit

- name: Publish Test Results
  uses: dorny/test-reporter@v1
  if: success() || failure()
  with:
    name: Go Tests
    path: junit-report.xml
    reporter: java-junit
```

#### Jenkins Pipeline Example
```groovy
stage('Test') {
    steps {
        sh 'make test-with-junit'
    }
    post {
        always {
            junit 'junit-report.xml'
        }
    }
}
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards
- Follow Go best practices and idioms
- Add comprehensive tests for new features with JUnit XML compatibility
- Update documentation for API changes
- Add Swagger annotations for new endpoints
- Follow the existing code structure
- Ensure all tests pass with `make test-with-junit`

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [HL7 FHIR](https://hl7.org/fhir/) for healthcare interoperability standards
- [Gin Web Framework](https://gin-gonic.com/) for HTTP routing
- [GORM](https://gorm.io/) for database ORM
- [Swagger](https://swagger.io/) for API documentation

---

## 📞 Support

For questions, issues, or contributions:
- **Issues:** [GitHub Issues](https://github.com/your-repo/go-fhir-demo/issues)
- **Discussions:** [GitHub Discussions](https://github.com/your-repo/go-fhir-demo/discussions)
- **Documentation:** [Wiki](https://github.com/your-repo/go-fhir-demo/wiki)

**Built with ❤️ for Healthcare Interoperability**

## Mocking

This application uses [github.com/uber-go/mock](https://github.com/uber-go/mock) for generating mocks in tests.

### Installing mockgen

To install the mock generator tool (`mockgen`), run:

```
go install go.uber.org/mock/mockgen@latest
```

Make sure your `$GOPATH/bin` (or `$GOBIN`) is in your system's `PATH` to use `mockgen` from anywhere.

### Installing gotestsum

To install gotestsum for enhanced test reporting and JUnit XML generation, run:

```
go install gotest.tools/gotestsum@latest
```

gotestsum provides:
- Enhanced test output formatting
- JUnit XML report generation for CI/CD integration
- Better test result visualization
- Support for various output formats
- Better test result visualization
- Support for various output formats
