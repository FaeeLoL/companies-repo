package database

import (
	"context"
	"fmt"

	"github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"

	"github.com/faeelol/companies-store/internal/app/database/migrations"
)

const (
	MigrDirectionUp   = "up"
	MigrDirectionDown = "down"
)

func MigrateDatabase(ctx context.Context, cfg *Config, direction string, logger logrus.FieldLogger) error {
	db, err := GetDB(cfg)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	var applied int

	switch direction {
	case MigrDirectionUp:
		applied, err = migrate.ExecContext(ctx, db.DB, "postgres", migrations.Migrations, migrate.Up)
	case MigrDirectionDown:
		applied, err = migrate.ExecContext(ctx, db.DB, "postgres", migrations.Migrations, migrate.Down)
	default:
		return fmt.Errorf("unknown migration direction: %s", direction)
	}

	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	logger.Infof("Applied %d migrations in direction '%s'\n", applied, direction)
	return nil
}
