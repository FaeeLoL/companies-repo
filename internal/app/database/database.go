// Package database provides functionality for managing the database connection, configuration, and queries.
package database

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
)

var (
	dbInstance *sqlx.DB
	once       sync.Once
)

func GetDB(cfg *Config) (*sqlx.DB, error) {
	var err error

	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
		)

		dbInstance, err = sqlx.Connect(cfg.Dialect, dsn)
		if err != nil {
			return
		}

		dbInstance.SetMaxOpenConns(25)
		dbInstance.SetMaxIdleConns(25)
		dbInstance.SetConnMaxLifetime(5 * 60)
	})

	return dbInstance, err
}
