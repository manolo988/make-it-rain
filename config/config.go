package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port            string        `mapstructure:"port"`
	Environment     string        `mapstructure:"environment"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxConnections  int           `mapstructure:"max_connections"`
	MinConnections  int           `mapstructure:"min_connections"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

type JWTConfig struct {
	SecretKey       string        `mapstructure:"secret_key"`
	ExpiryDuration  time.Duration `mapstructure:"expiry_duration"`
	RefreshDuration time.Duration `mapstructure:"refresh_duration"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	LogLevel    string `mapstructure:"log_level"`
	RateLimitRPS int    `mapstructure:"rate_limit_rps"`
}

var Cfg *Config

func LoadConfig(path string) error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.read_timeout", 10*time.Second)
	viper.SetDefault("server.write_timeout", 10*time.Second)
	viper.SetDefault("server.shutdown_timeout", 10*time.Second)

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_connections", 20)
	viper.SetDefault("database.min_connections", 2)
	viper.SetDefault("database.max_conn_lifetime", 1*time.Hour)
	viper.SetDefault("database.max_conn_idle_time", 30*time.Minute)

	viper.SetDefault("jwt.expiry_duration", 24*time.Hour)
	viper.SetDefault("jwt.refresh_duration", 7*24*time.Hour)

	viper.SetDefault("app.name", "Make It Rain API")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("app.rate_limit_rps", 100)

	viper.AutomaticEnv()

	// Explicitly bind environment variables
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.environment", "SERVER_ENVIRONMENT")
	viper.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	viper.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	viper.BindEnv("server.shutdown_timeout", "SERVER_SHUTDOWN_TIMEOUT")

	viper.BindEnv("database.host", "DATABASE_HOST")
	viper.BindEnv("database.port", "DATABASE_PORT")
	viper.BindEnv("database.user", "DATABASE_USER")
	viper.BindEnv("database.password", "DATABASE_PASSWORD")
	viper.BindEnv("database.name", "DATABASE_NAME")
	viper.BindEnv("database.ssl_mode", "DATABASE_SSL_MODE")
	viper.BindEnv("database.max_connections", "DATABASE_MAX_CONNECTIONS")
	viper.BindEnv("database.min_connections", "DATABASE_MIN_CONNECTIONS")
	viper.BindEnv("database.max_conn_lifetime", "DATABASE_MAX_CONN_LIFETIME")
	viper.BindEnv("database.max_conn_idle_time", "DATABASE_MAX_CONN_IDLE_TIME")

	viper.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")
	viper.BindEnv("jwt.expiry_duration", "JWT_EXPIRY_DURATION")
	viper.BindEnv("jwt.refresh_duration", "JWT_REFRESH_DURATION")

	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.version", "APP_VERSION")
	viper.BindEnv("app.log_level", "APP_LOG_LEVEL")
	viper.BindEnv("app.rate_limit_rps", "APP_RATE_LIMIT_RPS")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}