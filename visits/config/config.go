package config

import "time"

type (
	Config struct {
		APIHost         string        `env:"API_LISTEN_URL"       env-default:"0.0.0.0:8001"`
		ReadTimeout     time.Duration `env:"API_READ_TIMEOUT"     env-default:"5s"`
		WriteTimeout    time.Duration `env:"API_WRITE_TIMEOUT"    env-default:"5s"`
		ShutdownTimeout time.Duration `env:"API_SHUTDOWN_TIMEOUT" env-default:"5s"`
		IdleTimeout     time.Duration `env:"API_IDLE_TIMEOUT"     env-default:"120s"`

		Postgres Postgres
	}
	Postgres struct {
		Host         string        `env:"POSTGRES_HOST"              env-default:"0.0.0.0:5432"`
		Name         string        `env:"API_POSTGRES_DATABASE"      env-default:"postgres"`
		User         string        `env:"API_POSTGRES_USER"          env-default:"postgres"`
		Password     string        `env:"API_POSTGRES_PASSWORD"      env-default:"secret"`
		PoolSize     int           `env:"API_POSTGRES_POOL_SIZE"     env-default:"10"`
		MaxRetries   int           `env:"API_POSTGRES_MAX_RETRIES"   env-default:"5"`
		ReadTimeout  time.Duration `env:"API_POSTGRES_READ_TIMEOUT"  env-default:"10s"`
		WriteTimeout time.Duration `env:"API_POSTGRES_WRITE_TIMEOUT" env-default:"10s"`
		DisableTLS   bool          `env:"API_POSTGRES_DISABLE_TLS"   env-default:"true"`
	}
)
