package postgres

import (
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
	Host         string        `env:"POSTGRES_HOST"          env-default:"0.0.0.0:5432"`
	Name         string        `env:"POSTGRES_DATABASE"      env-default:"doctors"`
	User         string        `env:"POSTGRES_USER"          env-default:"postgres"`
	Password     string        `env:"POSTGRES_PASSWORD"      env-default:"secret"`
	PoolSize     int           `env:"POSTGRES_POOL_SIZE"     env-default:"10"`
	MaxRetries   int           `env:"POSTGRES_MAX_RETRIES"   env-default:"5"`
	ReadTimeout  time.Duration `env:"POSTGRES_READ_TIMEOUT"  env-default:"10s"`
	WriteTimeout time.Duration `env:"POSTGRES_WRITE_TIMEOUT" env-default:"10s"`
	DisableTLS   bool          `env:"POSTGRES_DISABLE_TLS"   env-default:"true"`
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
