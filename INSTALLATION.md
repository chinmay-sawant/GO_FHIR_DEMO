# Installation Guide

This document provides comprehensive installation instructions for the Go FHIR Demo Application, including Docker deployment, manual installation, and all required dependencies.

## 📋 Prerequisites

### System Requirements
- **Operating System**: Windows 10/11, macOS, or Linux
- **RAM**: Minimum 4GB, Recommended 8GB+
- **Disk Space**: At least 2GB free space
- **Network**: Internet connection for downloading dependencies

### Required Software

#### Core Requirements
- **Go 1.23+** - [Download from golang.org](https://golang.org/downloads/)
- **Git** - [Download from git-scm.com](https://git-scm.com/)
- **Docker Desktop** - [Download from docker.com](https://www.docker.com/products/docker-desktop/)

#### Database (Choose One)
- **Docker Compose** (Recommended) - Included with Docker Desktop
- **PostgreSQL 15+** - [Download from postgresql.org](https://www.postgresql.org/downloads/) (for manual installation)

#### Optional Tools
- **Make** - Build automation tool
  - **Windows**: Install via [Chocolatey](https://chocolatey.org/) (`choco install make`) or [Scoop](https://scoop.sh/) (`scoop install make`)
  - **macOS**: Install via Homebrew (`brew install make`)
  - **Linux**: Usually pre-installed (`sudo apt-get install make`)

## 🐳 Docker Installation (Recommended)

### 1. Clone the Repository
```bash
git clone <repository-url>
cd Go_FHIR_Demo
```

### 2. Environment Configuration
Create a `.env` file from the example:
```bash
# Copy the example environment file
cp .env.example .env
```

Edit the `.env` file with your configuration:
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
GIN_MODE=release

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

# Development Mode
DEV_MODE=false
```

### 3. Start All Services
```bash
# Start all services (PostgreSQL, Consul, Vault, Application)
docker-compose up -d

# View logs
docker-compose logs -f

# View logs for specific service
docker-compose logs -f fhir-api
```

### 4. Initialize Vault (Required)
After starting the services, Vault needs to be initialized:

```bash
# Initialize Vault (save the output!)
docker exec -it fhir_vault vault operator init

# Unseal Vault (use 3 different unseal keys from init output)
docker exec -it fhir_vault vault operator unseal <unseal_key_1>
docker exec -it fhir_vault vault operator unseal <unseal_key_2>
docker exec -it fhir_vault vault operator unseal <unseal_key_3>

# Login with root token
docker exec -it fhir_vault vault auth <root_token>

# Create a sample secret
docker exec -it fhir_vault vault kv put secret/myapp username=admin password=secret123
```

### 5. Verify Installation
Check that all services are running:
```bash
docker-compose ps
```

Access the application:
- **API**: [http://localhost:8080](http://localhost:8080)
- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **Health Check**: [http://localhost:8080/health](http://localhost:8080/health)
- **Consul UI**: [http://localhost:8500](http://localhost:8500)
- **Vault UI**: [http://localhost:8200](http://localhost:8200)

### 6. Stop Services
```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: This will delete all data)
docker-compose down -v
```

## 🔧 Manual Installation

### 1. Install Go
Download and install Go 1.23+ from [golang.org](https://golang.org/downloads/).

Verify installation:
```bash
go version
```

### 2. Install PostgreSQL
Download and install PostgreSQL 15+ from [postgresql.org](https://www.postgresql.org/downloads/).

Create database and user:
```sql
-- Connect to PostgreSQL as superuser
CREATE DATABASE fhir_demo;
CREATE USER fhir_user WITH PASSWORD 'fhir_password';
GRANT ALL PRIVILEGES ON DATABASE fhir_demo TO fhir_user;
GRANT ALL ON SCHEMA public TO fhir_user;
```

### 3. Install Consul (Optional)
Download Consul from [consul.io](https://www.consul.io/downloads):

```bash
# Start Consul in development mode
consul agent -dev -ui -client=0.0.0.0

# In another terminal, create a sample secret
consul kv put myapp/secret '{"username":"admin","password":"secret123"}'
```

### 4. Install Vault (Optional)
Download Vault from [vaultproject.io](https://www.vaultproject.io/downloads):

```bash
# Start Vault in development mode
vault server -dev

# In another terminal, set environment variables
export VAULT_ADDR='http://127.0.0.1:8200'
export VAULT_TOKEN="<dev_root_token_from_server_output>"

# Create a sample secret
vault kv put secret/myapp username=admin password=secret123
```

### 5. Install Migration Tool
Install golang-migrate for database migrations:

```bash
# Using Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Verify installation
migrate -version
```

### 6. Clone and Setup Application
```bash
# Clone repository
git clone <repository-url>
cd Go_FHIR_Demo

# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env
```

Edit `.env` with your local configuration:
```env
# Database Configuration (adjust for your local setup)
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

# Consul Configuration (if running locally)
CONSUL_ADDRESS=http://localhost:8500
CONSUL_KEY=myapp/secret

# Vault Configuration (if running locally)
VAULT_ADDRESS=http://localhost:8200
VAULT_TOKEN=<your_vault_token>
VAULT_SECRET_PATH=secret/data/myapp

# Development Mode
DEV_MODE=true
```

### 7. Run Database Migrations
```bash
# Run migrations
migrate -path migrations -database "postgres://fhir_user:fhir_password@localhost:5432/fhir_demo?sslmode=disable" up

# Verify migration status
migrate -path migrations -database "postgres://fhir_user:fhir_password@localhost:5432/fhir_demo?sslmode=disable" version
```

### 8. Start the Application
```bash
# Run the application
go run main.go

# Or build and run
go build -o fhir-demo
./fhir-demo
```

## 📦 Go Dependencies

The application uses the following Go packages (automatically installed with `go mod tidy`):

### Core Dependencies
```go
github.com/gin-gonic/gin v1.10.1                    // Web framework
github.com/spf13/viper v1.16.0                      // Configuration management
github.com/joho/godotenv v1.5.1                     // Environment variables
gorm.io/gorm v1.30.0                               // ORM
gorm.io/driver/postgres v1.5.7                     // PostgreSQL driver
github.com/sirupsen/logrus v1.9.3                  // Logging
```

### FHIR Support
```go
github.com/samply/golang-fhir-models/fhir-models v0.3.2  // FHIR R4 models
```

### API Documentation
```go
github.com/swaggo/gin-swagger v1.6.0                // Swagger integration
github.com/swaggo/files v1.0.1                      // Swagger files
github.com/swaggo/swag v1.16.4                      // Swagger generation
```

### Testing
```go
github.com/stretchr/testify v1.10.0                 // Testing framework
go.uber.org/mock v0.5.2                            // Mock generation
github.com/fergusstrange/embedded-postgres v1.31.0  // Embedded PostgreSQL for tests
```

## 🛠️ Development Tools Installation

### Install Testing Tools
```bash
# Install gotestsum for enhanced test reporting
go install gotest.tools/gotestsum@latest

# Install mockgen for generating mocks
go install go.uber.org/mock/mockgen@latest

# Install swag for generating Swagger documentation
go install github.com/swaggo/swag/cmd/swag@latest
```

### Install Migration Tools
```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Verify Tool Installation
```bash
# Verify all tools are installed
gotestsum --version
mockgen --version
swag --version
migrate -version
```

## 🔧 Makefile Commands

If you have `make` installed, you can use the following commands:

### Development Commands
```bash
make deps          # Download dependencies
make build         # Build the application
make run           # Run the application
make test          # Run tests
make clean         # Clean build artifacts
```

### Database Commands
```bash
make migrate-up    # Run database migrations up
make migrate-down  # Run database migrations down
make migrate-create name=your_migration_name  # Create new migration
```

### Testing Commands
```bash
make test-with-junit    # Run tests with JUnit XML output
make coverage          # Generate test coverage report
make coverage-detailed # Generate detailed coverage report
make mocks            # Generate all mocks
```

### Documentation Commands
```bash
make docs             # Generate Swagger documentation
make swagger          # Alias for docs
```

## 🐳 Docker Services Configuration

The `docker-compose.yml` includes the following services:

### PostgreSQL Database
- **Image**: `postgres:15-alpine`
- **Port**: `5432`
- **Volume**: `postgres_data`
- **Health Check**: Included

### Consul
- **Image**: `hashicorp/consul:1.21.0`
- **Port**: `8500`
- **Volume**: `consul_data`
- **UI**: Available at http://localhost:8500

### Vault
- **Image**: `hashicorp/vault:latest`
- **Port**: `8200`
- **Volume**: `vault_data`
- **Configuration**: Uses `vault/config/vault.hcl`
- **UI**: Available at http://localhost:8200

### Application
- **Build**: From Dockerfile
- **Port**: `8080`
- **Dependencies**: PostgreSQL, Consul, Vault
- **Environment**: Configured via `.env`

## 🔍 Troubleshooting

### Common Issues

#### Port Conflicts
If you encounter port conflicts:
```bash
# Check what's using the port
netstat -tulpn | grep :8080

# Stop the conflicting process or change the port in .env
SERVER_PORT=8081
```

#### Database Connection Issues
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Test connection manually
psql -h localhost -p 5432 -U fhir_user -d fhir_demo
```

#### Vault Initialization Issues
```bash
# Check Vault status
docker exec -it fhir_vault vault status

# Re-initialize if needed (WARNING: This will reset all data)
docker exec -it fhir_vault vault operator init -key-shares=3 -key-threshold=2
```

#### Permission Issues (Linux/macOS)
```bash
# Fix file permissions
chmod +x ./fhir-demo
sudo chown -R $USER:$USER ./logs
```

### Getting Help
If you encounter issues:
1. Check the application logs: `docker-compose logs fhir-api`
2. Verify all services are healthy: `docker-compose ps`
3. Check environment variables: `docker-compose config`
4. Review the troubleshooting section above

## 🧪 Testing the Installation

### Basic Health Check
```bash
# Test API health
curl http://localhost:8080/health

# Test Swagger UI
curl http://localhost:8080/swagger/index.html

# Test patient endpoint
curl http://localhost:8080/api/v1/patients
```

### Test External FHIR Integration
```bash
# Search external patients
curl "http://localhost:8080/api/v1/external-patients?family=Smith"

# Get external patient by ID
curl "http://localhost:8080/api/v1/external-patients/123"
```

### Test Consul Integration
```bash
# Get secret from Consul
curl http://localhost:8080/consul/secret
```

### Test Vault Integration (when implemented)
```bash
# Get secret from Vault
curl http://localhost:8080/vault/secret
```

## 📚 Next Steps

After successful installation:

1. **Explore the API**: Visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
2. **Create Patients**: Use the Swagger UI to create and manage FHIR patients
3. **Test External Integration**: Try searching external FHIR servers
4. **Configure Secrets**: Set up your secrets in Consul and Vault
5. **Customize Configuration**: Modify `config/config.json` for your environment
6. **Review Logs**: Check `logs/app.log` for application activity

## 🔐 Security Considerations

### Production Deployment
- Change default passwords and tokens
- Use proper SSL/TLS certificates
- Configure Vault authentication methods
- Set up proper firewall rules
- Use secrets management for sensitive data
- Enable audit logging

### Environment Variables
Never commit sensitive information to version control:
- Database passwords
- API keys
- Vault tokens
- External service credentials

Use environment variables or external secret management systems for production deployments.