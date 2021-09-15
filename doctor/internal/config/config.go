package config

import (
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/doctor/internal/repository/postgres"
)

//Config structure holds env variables as tags
type Config struct {
	APIHost               string        `env:"API_LISTEN_URL"               env-default:"0.0.0.0:8002"`
	DebugHost             string        `env:"API_DEBUG_URL"                env-default:"0.0.0.0:4002"`
	GRPCHost              string        `env:"GRPC_LISTEN_URL"              env-default:"0.0.0.0:6002"`
	UserGRPCClient        string        `env:"USERS_GRPC_URL"               env-default:"0.0.0.0:6000"`
	AppointmentGRPCClient string        `env:"APPOINTMENT_GRPC_URL"         env-default:"0.0.0.0:6001"`
	ReadTimeout           time.Duration `env:"API_READ_TIMEOUT"             env-default:"5s"`
	WriteTimeout          time.Duration `env:"API_WRITE_TIMEOUT"            env-default:"5s"`
	ShutdownTimeout       time.Duration `env:"API_SHUTDOWN_TIMEOUT"         env-default:"5s"`
	RedisURL              string        `env:"REDIS_URL"                    env-default:"localhost:6379"`
	RedisPassword         string        `env:"REDIS_PASSWORD"               env-default:""`

	Postgres postgres.Config
}
