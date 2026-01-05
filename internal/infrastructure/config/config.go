// Package config provides application configuration.
package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
	Telemetry TelemetryConfig
	Log       LogConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	Name            string        `mapstructure:"name"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	Debug           bool          `mapstructure:"debug"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret            string        `mapstructure:"secret"`
	Expiration        time.Duration `mapstructure:"expiration"`
	RefreshExpiration time.Duration `mapstructure:"refresh_expiration"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int           `mapstructure:"requests"`
	Window   time.Duration `mapstructure:"window"`
}

// TelemetryConfig holds TelemetryFlow configuration
type TelemetryConfig struct {
	APIKeyID       string `mapstructure:"api_key_id"`
	APIKeySecret   string `mapstructure:"api_key_secret"`
	Endpoint       string `mapstructure:"endpoint"`
	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load loads configuration from environment and config file
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")

	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.name", "orders")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	viper.SetDefault("database.debug", false)
	viper.SetDefault("jwt.expiration", "24h")
	viper.SetDefault("jwt.refresh_expiration", "168h")
	viper.SetDefault("ratelimit.requests", 100)
	viper.SetDefault("ratelimit.window", "1m")
	viper.SetDefault("telemetry.endpoint", "api.telemetryflow.id:4317")
	viper.SetDefault("telemetry.service_name", "Order-Service")
	viper.SetDefault("telemetry.service_version", "1.1.2")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// Bind environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("")

	// Environment variable mappings
	_ = viper.BindEnv("server.port", "SERVER_PORT")
	_ = viper.BindEnv("database.driver", "DB_DRIVER")
	_ = viper.BindEnv("database.host", "DB_HOST")
	_ = viper.BindEnv("database.port", "DB_PORT")
	_ = viper.BindEnv("database.name", "DB_NAME")
	_ = viper.BindEnv("database.user", "DB_USER")
	_ = viper.BindEnv("database.password", "DB_PASSWORD")
	_ = viper.BindEnv("database.ssl_mode", "DB_SSL_MODE")
	_ = viper.BindEnv("database.debug", "DB_DEBUG")
	_ = viper.BindEnv("jwt.secret", "JWT_SECRET")
	_ = viper.BindEnv("jwt.expiration", "JWT_EXPIRATION")
	_ = viper.BindEnv("telemetry.api_key_id", "TELEMETRYFLOW_API_KEY_ID")
	_ = viper.BindEnv("telemetry.api_key_secret", "TELEMETRYFLOW_API_KEY_SECRET")
	_ = viper.BindEnv("telemetry.endpoint", "TELEMETRYFLOW_ENDPOINT")
	_ = viper.BindEnv("telemetry.service_name", "TELEMETRYFLOW_SERVICE_NAME")

	_ = viper.BindEnv("log.level", "LOG_LEVEL")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
