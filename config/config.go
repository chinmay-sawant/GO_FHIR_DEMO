package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// JaegerConfig holds Jaeger-related configuration
type JaegerConfig struct {
	Endpoint    string `mapstructure:"JAEGER_ENDPOINT"`
	ServiceName string `mapstructure:"JAEGER_SERVICE_NAME"`
	Environment string `mapstructure:"JAEGER_ENVIRONMENT"`
	Enabled     bool   `mapstructure:"JAEGER_ENABLED"`
}

// KafkaConfig holds Kafka-related configuration
type KafkaConfig struct {
	Broker  string `mapstructure:"broker"`
	Topic   string `mapstructure:"topic"`
	GroupID string `mapstructure:"group_id"`
}

// Config holds all configuration for the application
type Config struct {
	DB     DatabaseConfig `mapstructure:"database"`
	Server ServerConfig   `mapstructure:"server"`
	Log    LogConfig      `mapstructure:"logging"`
	Consul ConsulConfig   `mapstructure:"consul"`
	Vault  VaultConfig    `mapstructure:"vault"`
	Redis  RedisConfig    `mapstructure:"redis"`
	Jaeger JaegerConfig   `mapstructure:"jaeger"`
	Kafka  KafkaConfig    `mapstructure:"kafka"`
	FHIR   FHIRConfig     `mapstructure:"fhir"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port                      string        `mapstructure:"port"`
	Mode                      string        `mapstructure:"mode"`
	ReadTimeout               time.Duration `mapstructure:"read_timeout"`
	WriteTimeout              time.Duration `mapstructure:"write_timeout"`
	ExternalFHIRServerBaseURL string        `mapstructure:"externalFHIRServerBaseURL"`
	DevMode                   bool          `mapstructure:"dev_mode"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// LogConfig holds logging-related configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	File   string `mapstructure:"file"`
}

// FHIRConfig holds FHIR-related configuration
type FHIRConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Version string `mapstructure:"version"`
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// ConsulConfig holds Consul-related configuration
type ConsulConfig struct {
	Address string `mapstructure:"address"`
	Key     string `mapstructure:"key"`
}

// VaultConfig holds Vault-related configuration
type VaultConfig struct {
	Address    string `mapstructure:"address"`
	Token      string `mapstructure:"token"`
	SecretPath string `mapstructure:"secret_path"`
}

// LoadConfig loads configuration from file and environment variables
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
	viper.SetDefault("jaeger.endpoint", "http://localhost:4318")
	viper.SetDefault("jaeger.service_name", "go-fhir-demo")
	viper.SetDefault("jaeger.environment", "development")
	viper.SetDefault("jaeger.enabled", true)
	viper.SetDefault("kafka.broker", "localhost:9092")
	viper.SetDefault("kafka.topic", "myapp-topic")
	viper.SetDefault("kafka.group_id", "myapp-group")

	// Bind environment variables
	_ = viper.BindEnv("server.port", "SERVER_PORT", "PORT")
	_ = viper.BindEnv("server.mode", "GIN_MODE")
	_ = viper.BindEnv("database.host", "DB_HOST")
	_ = viper.BindEnv("database.port", "DB_PORT")
	_ = viper.BindEnv("database.user", "DB_USER")
	_ = viper.BindEnv("database.password", "DB_PASSWORD")
	_ = viper.BindEnv("database.name", "DB_NAME")
	_ = viper.BindEnv("database.sslmode", "DB_SSLMODE")
	_ = viper.BindEnv("logging.level", "LOG_LEVEL")
	_ = viper.BindEnv("redis.host", "REDIS_HOST")
	_ = viper.BindEnv("redis.port", "REDIS_PORT")
	_ = viper.BindEnv("redis.password", "REDIS_PASSWORD")
	_ = viper.BindEnv("redis.db", "REDIS_DB")
	_ = viper.BindEnv("consul.address", "CONSUL_ADDRESS")
	_ = viper.BindEnv("consul.key", "CONSUL_KEY")
	_ = viper.BindEnv("vault.address", "VAULT_ADDRESS")
	_ = viper.BindEnv("vault.token", "VAULT_TOKEN")
	_ = viper.BindEnv("vault.secret_path", "VAULT_SECRET_PATH")
	_ = viper.BindEnv("jaeger.endpoint", "JAEGER_ENDPOINT")
	_ = viper.BindEnv("jaeger.service_name", "JAEGER_SERVICE_NAME")
	_ = viper.BindEnv("jaeger.environment", "JAEGER_ENVIRONMENT")
	_ = viper.BindEnv("jaeger.enabled", "JAEGER_ENABLED")
	_ = viper.BindEnv("kafka.broker", "KAFKA_BROKER")
	_ = viper.BindEnv("kafka.topic", "KAFKA_TOPIC")
	_ = viper.BindEnv("kafka.group_id", "KAFKA_GROUP_ID")

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
		config.DB.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		config.DB.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		config.DB.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.DB.Password = password
	}
	if name := os.Getenv("DB_NAME"); name != "" {
		config.DB.Name = name
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
	// Override with environment variables for Kafka
	if broker := os.Getenv("KAFKA_BROKER"); broker != "" {
		config.Kafka.Broker = broker
	}
	if topic := os.Getenv("KAFKA_TOPIC"); topic != "" {
		config.Kafka.Topic = topic
	}
	if groupID := os.Getenv("KAFKA_GROUP_ID"); groupID != "" {
		config.Kafka.GroupID = groupID
	}

	return &config, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}
