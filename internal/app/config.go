package app

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"

	"github.com/faeelol/companies-store/internal/app/database"
	"github.com/faeelol/companies-store/internal/app/kafka"
	"github.com/faeelol/companies-store/internal/app/rest"
)

type Config struct {
	Server *rest.Config
	DB     *database.Config
	Kafka  *kafka.Config
}

func NewConfig() *Config {
	return &Config{}
}

func LoadConfig(cfgPath string, cfg *Config) error {
	if err := loadGlobalConfig(cfgPath); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	cfg.Server = rest.LoadServerConfig()
	cfg.DB = database.LoadDatabaseConfig()
	cfg.Kafka = kafka.LoadKafkaConfig()

	return nil
}

func loadGlobalConfig(configPath string) error {
	dir := filepath.Dir(configPath)
	file := filepath.Base(configPath)

	filename := strings.TrimSuffix(file, filepath.Ext(file))
	ext := filepath.Ext(file)[1:]

	// setup Viper
	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(dir)

	// read config
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	return nil
}
