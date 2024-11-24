// Package app initializes and manages the core application components.
package app

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/faeelol/companies-store/internal/app/database"
	"github.com/faeelol/companies-store/internal/app/kafka"
	"github.com/faeelol/companies-store/internal/app/rest"
)

func LoadConfigInitLoggerAndDo(configPath string, cfg *Config, f func(logger *logrus.Logger) error) error {
	err := LoadConfig(configPath, cfg)
	if err != nil {
		return err
	}

	logger := InitLogger()

	return f(logger)
}

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return logger
}

func StartHTTPService(ctx context.Context, cfg *Config, logger *logrus.Logger) error {
	db, err := database.GetDB(cfg.DB)
	if err != nil {
		return err
	}

	kafkaProducer := kafka.NewProducer(cfg.Kafka)
	defer func(kProducer *kafka.Producer) {
		_ = kProducer.Close()
	}(kafkaProducer)

	server := rest.NewServer(cfg.Server, logger, db, kafkaProducer)

	return server.Start(ctx, logger)
}

func MigrateDatabase(ctx context.Context, cfg *Config, direction string, logger *logrus.Logger) error {
	return database.MigrateDatabase(ctx, cfg.DB, direction, logger)
}
