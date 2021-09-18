// Package postgres provides ...
package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	// migration postgres drivers.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateUp runs migration and applies everything new to the DB provided in dsn string
func MigrateUp(migrationsPath, dsn string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn)
	if err != nil {
		return fmt.Errorf("migration failed, %v", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			return fmt.Errorf("migration failed, %v", err)
		}
	}
	return nil
}
