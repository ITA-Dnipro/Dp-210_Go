package config

type Env struct {
	Connection string `env:"DB_CONNECTION" env-default:"postgres://postgres:secret@0.0.0.0:5432/test?sslmode=disable&timezone=utc"`
	SqlDriver  string `env:"SQL_DRIVER" env-default:"pgx"`
	Host       string `env:"HOST" env-default:"localhost"`
	Port       string `env:"PORT" env-default:"8000"`
}
