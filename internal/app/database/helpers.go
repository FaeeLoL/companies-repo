package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TxFunc func(tx *sqlx.Tx) error

// WithinTransaction runs the provided function within a transaction.
func WithinTransaction(ctx context.Context, db *sqlx.DB, fn TxFunc) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
