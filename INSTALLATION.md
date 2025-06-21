# Installation Guide

This guide provides detailed instructions for setting up and running the Go FHIR Demo Application.

## 📋 Prerequisites

### Required Software

#### 1. **Go 1.23+**
- **Download:** [https://golang.org/downloads/](https://golang.org/downloads/)
- **Verify installation:**
  ```bash
  go version
  ```

#### 2. **Docker & Docker Compose**
- **Docker Desktop (Windows/Mac):** [https://docs.docker.com/desktop/](https://docs.docker.com/desktop/)
- **Docker Engine (Linux):** [https://docs.docker.com/engine/install/](https://docs.docker.com/engine/install/)
- **Verify installation:**
  ```bash
  docker --version
  docker-compose --version
  ```

#### 3. **Git**
- **Download:** [https://git-scm.com/downloads](https://git-scm.com/downloads)
- **Verify installation:**
  ```bash
  git --version
  ```

#### 4. **PostgreSQL 15+** (for manual setup)
- **Download:** [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
- **Verify installation:**
  ```bash
  psql --version
  ```

### Optional Development Tools

#### 1. **golang-migrate** (Database migrations)
```bash
# Using Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Or download binary from: https://github.com/golang-migrate/migrate/releases
```

#### 2. **swag** (Swagger documentation generator)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

#### 3. **gotestsum** (Enhanced test output)
```bash
go install gotest.tools/gotestsum@latest
```

#### 4. **gomock** (Mock generation for testing)
```bash
go install go.uber.org/mock/mockgen@latest
```

#### 5. **Make** (Windows - optional)
- **Chocolatey:** `choco install make`
- **Or use batch scripts provided in the project**

## 🚀 Installation Methods

### Method 1: Docker Compose (Recommended)

This method sets up all services including PostgreSQL, Consul, and Vault automatically.

#### Step 1: Clone Repository
```bash
git clone <repository-url>
cd Go_FHIR_Demo
```

#### Step 2: Environment Configuration
Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` file with your configurations:
```env
# Database Configuration
DB_HOST=postgres
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

# Consul Configuration
CONSUL_ADDRESS=http://consul:8500
CONSUL_KEY=myapp/secret

# Vault Configuration
VAULT_ADDRESS=http://vault:8200
VAULT_TOKEN=root
VAULT_SECRET_PATH=secret/data/myapp
```

#### Step 3: Start Services
```bash
# Start all services
docker-compose up -d

# Or start specific services
docker-compose up -d postgres consul vault

# Build and start the API
docker-compose up --build fhir-api
```

#### Step 4: Initialize Vault (First Time Only)
```bash
# Initialize Vault
docker exec -it fhir_vault vault operator init

# Save the unseal keys and root token!
# Unseal Vault (repeat 3 times with different keys)
docker exec -it fhir_vault vault operator unseal <unseal_key_1>
docker exec -it fhir_vault vault operator unseal <unseal_key_2>
docker exec -it fhir_vault vault operator unseal <unseal_key_3>
```

#### Step 5: Verify Installation
- **API:** http://localhost:8080/health
- **Swagger UI:** http://localhost:8080/swagger/index.html
- **Consul UI:** http://localhost:8500
- **Vault UI:** http://localhost:8200/ui

### Method 2: Manual Setup

For development or when you want more control over individual components.

#### Step 1: Clone Repository
```bash
git clone <repository-url>
cd Go_FHIR_Demo
```

#### Step 2: Install Go Dependencies
```bash
go mod tidy
```

#### Step 3: Database Setup

##### Option A: Using Docker for PostgreSQL only
```bash
docker run --name fhir_postgres \
  -e POSTGRES_DB=fhir_demo \
  -e POSTGRES_USER=fhir_user \
  -e POSTGRES_PASSWORD=fhir_password \
  -p 5432:5432 \
  -d postgres:15-alpine
```

##### Option B: Manual PostgreSQL Setup
1. Install PostgreSQL
2. Create database and user:
```sql
CREATE DATABASE fhir_demo;
CREATE USER fhir_user WITH PASSWORD 'fhir_password';
GRANT ALL PRIVILEGES ON DATABASE fhir_demo TO fhir_user;
```

#### Step 4: Environment Configuration
Create `.env` file:
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

# Consul Configuration (optional)
CONSUL_ADDRESS=http://localhost:8500
CONSUL_KEY=myapp/secret

# Vault Configuration (optional)
VAULT_ADDRESS=http://localhost:8200
VAULT_TOKEN=root
VAULT_SECRET_PATH=secret/data/myapp
```

#### Step 5: Run Database Migrations
```bash
migrate -path migrations -database "postgres://fhir_user:fhir_password@localhost:5432/fhir_demo?sslmode=disable" up
```

#### Step 6: Generate Swagger Documentation
```bash
swag init --parseDependency --parseDepth 99
```

#### Step 7: Run Application
```bash
go run main.go
```

### Method 3: Development Setup with Hot Reload

#### Step 1: Install Air (Hot Reload Tool)
```bash
go install github.com/cosmtrek/air@latest
```

#### Step 2: Follow Manual Setup (Steps 1-6)

#### Step 3: Run with Hot Reload
```bash
air
```

## 🔧 Configuration

### Configuration Files

#### 1. **config/config.json** (Default values)
```json
{
  "server": {
    "port": "8080",
    "mode": "debug",
    "read_timeout": "10s",
    "write_timeout": "10s"
  },
  "database": {
    "host": "localhost",
    "port": "5432",
    "sslmode": "disable",
    "max_idle_conns": 10,
    "max_open_conns": 100,
    "conn_max_lifetime": "1h"
  },
  "consul": {
    "address": "http://localhost:8500",
    "key": "myapp/secret"
  },
  "vault": {
    "address": "http://localhost:8200",
    "token": "root",
    "secret_path": "secret/data/myapp"
  }
}
```

#### 2. **vault/config/vault.hcl** (Vault Configuration)
```hcl
storage "file" {
  path = "/vault/file"
}

listener "tcp" {
  address     = "0.0.0.0:8200"
  tls_disable = 1
}

disable_mlock = true
ui = true
api_addr = "http://fhir_vault:8200"
cluster_addr = "http://fhir_vault:8201"
```

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Database host | `localhost` | Yes |
| `DB_PORT` | Database port | `5432` | Yes |
| `DB_USER` | Database username | - | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | - | Yes |
| `DB_SSLMODE` | SSL mode | `disable` | No |
| `SERVER_PORT` | HTTP server port | `8080` | No |
| `GIN_MODE` | Gin mode (`debug`/`release`) | `debug` | No |
| `LOG_LEVEL` | Log level | `info` | No |
| `EXTERNAL_FHIR_SERVER_BASE_URL` | External FHIR server URL | - | Yes |
| `CONSUL_ADDRESS` | Consul server address | `http://localhost:8500` | No |
| `CONSUL_KEY` | Consul KV key | `myapp/secret` | No |
| `VAULT_ADDRESS` | Vault server address | `http://localhost:8200` | No |
| `VAULT_TOKEN` | Vault authentication token | `root` | No |
| `VAULT_SECRET_PATH` | Vault secret path | `secret/data/myapp` | No |

## 🧪 Testing Setup

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed output
gotestsum --format testname

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Generate Test Reports
```bash
# Generate JUnit XML report
gotestsum --junitfile junit-report.xml

# Generate coverage report
go test -coverprofile=coverage.out ./...
```

## 🛠️ Development Commands

### Using Makefile
```bash
# Install dependencies
make deps

# Build application
make build

# Run application
make run

# Run tests
make test

# Generate mocks
make generate-mocks

# Generate swagger docs
make swagger

# Clean build artifacts
make clean
```

### Manual Commands
```bash
# Install dependencies
go mod tidy

# Build application
go build -o bin/fhir-demo main.go

# Run application
go run main.go

# Run tests
go test ./...

# Generate mocks
mockgen -source=internal/domain/patient.go -destination=internal/domain/mocks/mock_patient.go

# Generate swagger docs
swag init --parseDependency --parseDepth 99
```

## 🔍 Troubleshooting

### Common Issues

#### 1. **Port Already in Use**
```bash
# Check what's using the port
netstat -ano | findstr :8080
netstat -ano | findstr :8200
netstat -ano | findstr :5432

# Kill the process
taskkill /PID <PID> /F
```

#### 2. **Database Connection Issues**
- Verify PostgreSQL is running
- Check database credentials
- Ensure database exists
- Check firewall settings

#### 3. **Vault Initialization Issues**
```bash
# Check Vault logs
docker logs fhir_vault

# Restart Vault container
docker-compose restart vault

# Initialize Vault manually
docker exec -it fhir_vault vault operator init
```

#### 4. **Migration Issues**
```bash
# Check migration status
migrate -path migrations -database "postgres://..." version

# Force migration version
migrate -path migrations -database "postgres://..." force <version>
```

### Docker Issues

#### 1. **Container Won't Start**
```bash
# Check container logs
docker logs fhir_api

# Rebuild containers
docker-compose up --build

# Remove containers and volumes
docker-compose down -v
docker system prune -f
```

#### 2. **Volume Issues**
```bash
# Remove all volumes
docker-compose down -v

# List volumes
docker volume ls

# Remove specific volume
docker volume rm go_fhir_demo_postgres_data
```

## 📚 Additional Resources

- **FHIR R4 Specification:** [https://hl7.org/fhir/R4/](https://hl7.org/fhir/R4/)
- **Gin Framework:** [https://gin-gonic.com/](https://gin-gonic.com/)
- **GORM Documentation:** [https://gorm.io/docs/](https://gorm.io/docs/)
- **Docker Compose:** [https://docs.docker.com/compose/](https://docs.docker.com/compose/)
- **Consul Documentation:** [https://www.consul.io/docs](https://www.consul.io/docs)
- **Vault Documentation:** [https://www.vaultproject.io/docs](https://www.vaultproject.io/docs)

## 🆘 Getting Help

If you encounter issues:

1. Check the logs: `docker-compose logs fhir-api`
2. Verify all services are running: `docker-compose ps`
3. Check the troubleshooting section above
4. Review environment variables and configuration
5. Ensure all prerequisites are installed correctly