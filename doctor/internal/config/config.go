package config

import (
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/repository/postgres"
)

type Config struct {
	APIHost         string        `env:"API_LISTEN_URL"       env-default:"0.0.0.0:8000"`
	DebugHost       string        `env:"API_DEBUG_URL"        env-default:"0.0.0.0:4000"`
	ReadTimeout     time.Duration `env:"API_READ_TIMEOUT"     env-default:"5s"`
	WriteTimeout    time.Duration `env:"API_WRITE_TIMEOUT"    env-default:"5s"`
	ShutdownTimeout time.Duration `env:"API_SHUTDOWN_TIMEOUT" env-default:"5s"`

	Postgres postgres.Config
}
