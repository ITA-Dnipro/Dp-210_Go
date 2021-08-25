// Package postgres provides ...
package postgres

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/ITA-Dnipro/Dp-210_Go/visits/config"
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

func Open(cfg config.Postgres) (*sql.DB, error) {
	q := make(url.Values)
	if cfg.DisableTLS {
		q.Set("sslmode", "disable")
	}
	q.Set("timezone", "utc")

	// nolint:exhaustivestruct
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sql.Open("pgx", u.String())
	if err != nil {
		return nil, fmt.Errorf("creating db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()

		return nil, fmt.Errorf("ping db %s : %w", u.String(), err)
	}

	return db, nil
}
