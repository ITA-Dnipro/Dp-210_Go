package postgres

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/ITA-Dnipro/Dp-210_Go/visits/config"

	// The database driver in use.
	_ "github.com/jackc/pgx/v4/stdlib"
)

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
