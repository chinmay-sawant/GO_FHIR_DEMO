package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
	FHIR     FHIRConfig     `json:"fhir"`
	Redis    RedisConfig    `json:"redis"`
	Consul   ConsulConfig   `json:"consul"`
	Vault    VaultConfig    `json:"vault"`
	Jaeger   JaegerConfig   `json:"jaeger"`
}

type ServerConfig struct {
	Port                      string        `json:"port"`
	Mode                      string        `json:"mode"`
	ReadTimeout               time.Duration `json:"read_timeout"`
	WriteTimeout              time.Duration `json:"write_timeout"`
	ExternalFHIRServerBaseURL string        `json:"external_fhir_server_base_url"`
	DevMode                   bool          `json:"dev_mode"`
}

type DatabaseConfig struct {
	Host            string        `json:"host"`
	Port            string        `json:"port"`
	User            string        `json:"user"`
	Password        string        `json:"password"`
	Name            string        `json:"name"`
	SSLMode         string        `json:"sslmode"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	File   string `json:"file"`
}

type FHIRConfig struct {
	BaseURL string `json:"base_url"`
	Version string `json:"version"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type ConsulConfig struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

type VaultConfig struct {
	Address    string `json:"address"`
	Token      string `json:"token"`
	SecretPath string `json:"secret_path"`
}

type JaegerConfig struct {
	Endpoint    string `json:"endpoint"`
	ServiceName string `json:"service_name"`
	Environment string `json:"environment"`
	Enabled     bool   `json:"enabled"`
}

func Load() (*Config, error) {
	// Load .env file from the root directory if it exists
	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", "1h")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.file", "logs/app.log")
	viper.SetDefault("fhir.base_url", "/api/v1")
	viper.SetDefault("fhir.version", "R4")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("consul.address", "http://localhost:8500")
	viper.SetDefault("consul.key", "myapp/secret")
	viper.SetDefault("vault.address", "http://localhost:8200")
	viper.SetDefault("vault.token", "root")
	viper.SetDefault("vault.secret_path", "secret/data/myapp")
	viper.SetDefault("jaeger.endpoint", "http://localhost:14268/api/traces")
	viper.SetDefault("jaeger.service_name", "go-fhir-demo")
	viper.SetDefault("jaeger.environment", "development")
	viper.SetDefault("jaeger.enabled", true)

	// Bind environment variables
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "GIN_MODE")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
	viper.BindEnv("logging.level", "LOG_LEVEL")
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")
	viper.BindEnv("consul.address", "CONSUL_ADDRESS")
	viper.BindEnv("consul.key", "CONSUL_KEY")
	viper.BindEnv("vault.address", "VAULT_ADDRESS")
	viper.BindEnv("vault.token", "VAULT_TOKEN")
	viper.BindEnv("vault.secret_path", "VAULT_SECRET_PATH")
	viper.BindEnv("jaeger.endpoint", "JAEGER_ENDPOINT")
	viper.BindEnv("jaeger.service_name", "JAEGER_SERVICE_NAME")
	viper.BindEnv("jaeger.environment", "JAEGER_ENVIRONMENT")
	viper.BindEnv("jaeger.enabled", "JAEGER_ENABLED")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Override with environment variables for database
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.Database.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.Database.Name = name
	}
	// Override with environment variables for Redis
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		config.Redis.Port = port
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}
	// Override with environment variables for Consul
	if addr := os.Getenv("CONSUL_ADDRESS"); addr != "" {
		config.Consul.Address = addr
	}
	if key := os.Getenv("CONSUL_KEY"); key != "" {
		config.Consul.Key = key
	}
	// Override with environment variables for Vault
	if addr := os.Getenv("VAULT_ADDRESS"); addr != "" {
		config.Vault.Address = addr
	}
	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		config.Vault.Token = token
	}
	if path := os.Getenv("VAULT_SECRET_PATH"); path != "" {
		config.Vault.SecretPath = path
	}
	if key := os.Getenv("DEV_MODE"); key != "" {
		if key == "true" {
			config.Server.DevMode = true
		}
	}
	// Override with environment variables for Jaeger
	if endpoint := os.Getenv("JAEGER_ENDPOINT"); endpoint != "" {
		config.Jaeger.Endpoint = endpoint
	}
	if serviceName := os.Getenv("JAEGER_SERVICE_NAME"); serviceName != "" {
		config.Jaeger.ServiceName = serviceName
	}
	if environment := os.Getenv("JAEGER_ENVIRONMENT"); environment != "" {
		config.Jaeger.Environment = environment
	}
	if enabled := os.Getenv("JAEGER_ENABLED"); enabled == "false" {
		config.Jaeger.Enabled = false
	}

	return &config, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
