package migrations

import "github.com/rubenv/sql-migrate"

func NewMigration1732452571InitialMigration() *migrate.Migration {
	return &migrate.Migration{
		Id: "1732452571_initial_migration.go",
		Up: []string{
			`
				CREATE TYPE company_type AS ENUM (
				'Corporations',
				'NonProfit',
				'Cooperative',
				'Sole Proprietorship'
			);
			
			CREATE TABLE companies (
				id UUID PRIMARY KEY,
				name VARCHAR(15) NOT NULL UNIQUE,
				description TEXT,
				employees_count INT NOT NULL,
				registered BOOLEAN NOT NULL,
				type company_type NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
				`,
		},
		Down: []string{
			`
			DROP TABLE IF EXISTS companies;
			
			DROP TYPE IF EXISTS company_type;
			`,
		},
	}
}
