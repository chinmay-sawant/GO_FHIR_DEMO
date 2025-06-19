# Go FHIR Demo Application

A comprehensive Go Gin framework application with FHIR (Fast Healthcare Interoperability Resources) R4 support, featuring a PostgreSQL database, automatic API documentation with Swagger, **external FHIR server integration**, and production-ready architecture.

## 🚀 Features

### Core Features
- **RESTful API** for Patient resources (GET, POST, PUT, DELETE, PATCH)
- **FHIR R4 Compliance** with custom FHIR data structures using standard JSON
- **External FHIR Server Integration** - Connect to and query external FHIR servers (like HAPI FHIR)
- **FHIR Client Package** - Reusable client for external FHIR server communication
- **PostgreSQL Database** with GORM ORM for robust data persistence
- **Database Migrations** using [golang-migrate](https://github.com/golang-migrate/migrate)
- **Swagger/OpenAPI Documentation** with interactive UI
- **Automatic Data Seeding** with sample FHIR patient records
- **Structured Logging** with Logrus
- **Configuration Management** with Viper and environment variables
- **Request/Response Middleware** for performance monitoring
- **Clean Architecture** with proper separation of concerns

### Testing & Quality Assurance
- **Comprehensive Test Suite** with unit and integration tests
- **JUnit XML Reports** using [gotestsum](https://github.com/gotestyourself/gotestsum) for CI/CD integration
- **Test Coverage Reports** with HTML and console output
- **Mock Generation** using [uber-go/mock](https://github.com/uber-go/mock) for unit testing
- **Automated Test Reporting** compatible with Jenkins, GitHub Actions, and other CI tools

### Technical Features
- Docker support with docker-compose
- Makefile for common development tasks
- Environment-based configuration with `.env` file support
- JSONB storage for efficient FHIR data querying
- HTTP client with timeout and error handling for external FHIR servers
- Comprehensive error handling and validation
- Production-ready logging and monitoring

## 🏗️ Architecture & Project Structure

```
├── config/                    # Configuration management
│   ├── config.go             # Configuration loader with Viper
│   └── config.json           # Default configuration
├── docs/                     # Auto-generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── examples/                 # Sample FHIR data
│   └── sample_patient.json   # Example patient resource
├── internal/                 # Private application code
│   ├── api/
│   │   ├── handlers/         # HTTP request handlers
│   │   │   ├── patient_handler.go
│   │   │   └── external_patient_handler.go  # External FHIR server handlers
│   │   └── routes/           # Route definitions
│   │       └── routes.go
│   ├── domain/               # Domain models and business entities
│   │   ├── patient.go        # FHIR Patient domain model
│   │   └── external_patient.go  # External patient service interface
│   ├── middleware/           # HTTP middleware
│   │   └── middleware.go     # Logging and timing middleware
│   ├── repository/           # Data access layer
│   │   └── patient_repository.go
│   └── service/              # Business logic layer
│       ├── patient_service.go
│       └── external_patient_service.go  # External FHIR server service
├── logs/                     # Application logs
├── migrations/               # Database schema migrations
│   ├── 000001_create_patients_table.up.sql
│   └── 000001_create_patients_table.down.sql
├── pkg/                      # Shared/reusable packages
│   ├── database/             # Database connection utilities
│   ├── fhirclient/           # FHIR client for external servers
│   │   └── client.go         # HTTP client for FHIR R4 servers
│   ├── logger/               # Logging utilities
│   └── utils/                # Common utility functions
├── docker-compose.yml        # Docker services definition
├── Dockerfile               # Container build instructions
├── Makefile                 # Development automation
└── main.go                  # Application entry point
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

## 📋 Prerequisites

Before you begin, ensure you have the following installed on your Windows system:

- **Go 1.21+** - [Download from golang.org](https://golang.org/downloads/)
- **PostgreSQL 12+** - [Download from postgresql.org](https://www.postgresql.org/downloads/)
- **Git** - [Download from git-scm.com](https://git-scm.com/)
- **Make** (optional) - [Install via Chocolatey](https://chocolatey.org/) or use the provided batch scripts

## 🚀 Quick Start

### 1. Clone the Repository
```cmd
git clone <repository-url>
cd Go_FHIR_Demo
```

### 2. Environment Setup

Copy the example environment file and configure your settings:
```cmd
copy .env.example .env
```

Edit the `.env` file with your database credentials and external FHIR server URL:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=fhir_user
DB_PASSWORD=fhir_password
DB_NAME=fhir_demo
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# External FHIR Server Configuration
EXTERNAL_FHIR_SERVER_BASE_URL=http://hapi.fhir.org/baseR4

# Logging Configuration
LOG_LEVEL=info
```

### 3. Database Setup

#### Option A: Using Docker (Recommended)
```cmd
docker-compose up -d postgres
```

#### Option B: Manual PostgreSQL Setup
1. Create a PostgreSQL database named `fhir_demo`
2. Create a user with appropriate permissions
```sql
CREATE DATABASE fhir_demo;
CREATE USER fhir_user WITH PASSWORD 'fhir_password';
GRANT ALL PRIVILEGES ON DATABASE fhir_demo TO fhir_user;
```

### 4. Install Dependencies
```cmd
go mod tidy
```

### 5. Install Migration Tool

**Option A: Using Go Install (Recommended)**
```cmd
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Option B: Download Binary (Windows)**
1. Download the latest Windows release from: https://github.com/golang-migrate/migrate/releases
2. Extract `migrate.exe` and place it in a folder (e.g., `C:\tools\migrate.exe`)
3. Add that folder to your system PATH
4. Verify installation:
```cmd
migrate --version
```

### 6. Run Database Migrations
```cmd
migrate -path migrations -database "postgres://fhir_user:fhir_password@localhost:5432/fhir_demo?sslmode=disable" up
```

### 7. Run the Application
```cmd
go run main.go
```

The application will start on `http://localhost:8080` and automatically seed sample patient data.

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

### Local Patient Resource Endpoints

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `GET` | `/api/v1/patients` | Get all patients with pagination | - |
| `GET` | `/api/v1/patients/{id}` | Get patient by ID | - |
| `POST` | `/api/v1/patients` | Create new patient | FHIR Patient JSON |
| `PUT` | `/api/v1/patients/{id}` | Update entire patient resource | FHIR Patient JSON |
| `PATCH` | `/api/v1/patients/{id}` | Partially update patient | Partial FHIR Patient JSON |
| `DELETE` | `/api/v1/patients/{id}` | Delete patient (soft delete) | - |

### External FHIR Server Endpoints

| Method | Endpoint | Description | Query Parameters |
|--------|----------|-------------|------------------|
| `GET` | `/api/v1/external-patients/{id}` | Get patient from external FHIR server by ID | - |
| `GET` | `/api/v1/external-patients` | Search patients on external FHIR server | `name`, `family`, `given`, `birthdate`, `gender` |

### Health Check Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Application health status |
| `GET` | `/swagger/index.html` | Interactive API documentation |

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

#### Search External FHIR Server
```bash
curl -X GET "http://localhost:8080/api/v1/external-patients?family=Smith&gender=female"
```

#### Get All Local Patients
```bash
curl -X GET http://localhost:8080/api/v1/patients
```

#### Get Patient by ID
```bash
curl -X GET http://localhost:8080/api/v1/patients/1
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

### Configuration File
The application also supports JSON configuration via `config/config.json` for default values. Environment variables take precedence over configuration file settings.

#### External FHIR Server Configuration
Configure the base URL for the external FHIR server in your environment:
```env
EXTERNAL_FHIR_SERVER_BASE_URL=http://hapi.fhir.org/baseR4
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
# Start all services (PostgreSQL + Application)
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
7. **Documentation** - Add Swagger annotations

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
