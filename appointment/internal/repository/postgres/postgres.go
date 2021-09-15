// Package postgres provides ...
package postgres

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"

	// The database driver in use.
	_ "github.com/jackc/pgx/v4/stdlib"
	// migration postgres drivers.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	Host         string        `envconfig:"POSTGRES_HOST"          default:"0.0.0.0:5432"`
	Name         string        `envconfig:"POSTGRES_DATABASE"      default:"appointments"`
	User         string        `envconfig:"POSTGRES_USER"          default:"postgres"`
	Password     string        `envconfig:"POSTGRES_PASSWORD"      default:"secret"`
	PoolSize     int           `envconfig:"POSTGRES_POOL_SIZE"     default:"10"`
	MaxRetries   int           `envconfig:"POSTGRES_MAX_RETRIES"   default:"5"`
	ReadTimeout  time.Duration `envconfig:"POSTGRES_READ_TIMEOUT"  default:"10s"`
	WriteTimeout time.Duration `envconfig:"POSTGRES_WRITE_TIMEOUT" default:"10s"`
	DisableTLS   bool          `envconfig:"POSTGRES_DISABLE_TLS"   default:"true"`
}

func (cfg *Config) String() string {
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
	return u.String()
}

func Open(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		return nil, fmt.Errorf("creating db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db %s : %w", cfg.String(), err)
	}
	return db, nil
}

// MigrateUp runs migration and applies everything new to the DB provided in dsn string
func MigrateUp(migrationsPath string, cfg Config) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		cfg.String(),
	)
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
