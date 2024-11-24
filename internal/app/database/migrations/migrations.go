// Package migrations contains database queries for initializing and updating the database schema.
package migrations

import "github.com/rubenv/sql-migrate"

var Migrations = &migrate.MemoryMigrationSource{
	Migrations: []*migrate.Migration{
		NewMigration1732452571InitialMigration(),
	},
}
