package config

import (
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/repository/postgres"
)

type Config struct {
	APIHost         string        `envconfig:"API_LISTEN_URL"       default:"0.0.0.0:8001"`
	GRPCHost        string        `envconfig:"GRPC_LISTEN_URL"      default:"0.0.0.0:6000"`
	DebugHost       string        `envconfig:"API_DEBUG_URL"        default:"0.0.0.0:4000"`
	DocrotGRPCHost  string        `envconfig:"DOCTORS_GRPC_URL"     default:"0.0.0.0:6002"`
	KafkaBrokers    []string      `envconfig:"KAFKA_BROKERS"        default:"0.0.0.0:9091"`
	ReadTimeout     time.Duration `envconfig:"API_READ_TIMEOUT"     default:"5s"`
	WriteTimeout    time.Duration `envconfig:"API_WRITE_TIMEOUT"    default:"5s"`
	ShutdownTimeout time.Duration `envconfig:"API_SHUTDOWN_TIMEOUT" default:"5s"`

	Postgres postgres.Config
}
