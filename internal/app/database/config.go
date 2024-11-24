package database

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Dialect       string
	Host          string
	Port          int
	User          string
	Password      string
	Database      string
	SSLMode       string
	Timeout       time.Duration
	MigrationsDir string
}

func NewConfig() *Config {
	return &Config{}
}

// LoadDatabaseConfig uploads the database configuration from the configuration file.
func LoadDatabaseConfig() *Config {
	viper.SetDefault("database.dialect", "postgres")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "user")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "companies-store")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.timeout", "5s")
	viper.SetDefault("database.migrations_dir", "./internal/app/database/migrations")

	return &Config{
		Dialect:       viper.GetString("database.dialect"),
		Host:          viper.GetString("database.host"),
		Port:          viper.GetInt("database.port"),
		User:          viper.GetString("database.user"),
		Password:      viper.GetString("database.password"),
		Database:      viper.GetString("database.name"),
		SSLMode:       viper.GetString("database.sslmode"),
		Timeout:       viper.GetDuration("database.timeout"),
		MigrationsDir: viper.GetString("database.migrations_dir"),
	}
}
