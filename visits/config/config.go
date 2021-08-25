package config

import (
	"fmt"
	"time"
)

type (
	Config struct {
		APIHost         string        `env:"API_LISTEN_URL"       env-default:"0.0.0.0:8001"`
		ReadTimeout     time.Duration `env:"API_READ_TIMEOUT"     env-default:"5s"`
		WriteTimeout    time.Duration `env:"API_WRITE_TIMEOUT"    env-default:"5s"`
		ShutdownTimeout time.Duration `env:"API_SHUTDOWN_TIMEOUT" env-default:"5s"`
		IdleTimeout     time.Duration `env:"API_IDLE_TIMEOUT"     env-default:"120s"`

		Postgres Postgres
		Kaffka   Kaffka
	}
	Postgres struct {
		Host         string        `env:"POSTGRES_HOST"              env-default:"0.0.0.0:5432"`
		Name         string        `env:"API_POSTGRES_DATABASE"      env-default:"test"`
		User         string        `env:"API_POSTGRES_USER"          env-default:"postgres"`
		Password     string        `env:"API_POSTGRES_PASSWORD"      env-default:"secret"`
		PoolSize     int           `env:"API_POSTGRES_POOL_SIZE"     env-default:"10"`
		MaxRetries   int           `env:"API_POSTGRES_MAX_RETRIES"   env-default:"5"`
		ReadTimeout  time.Duration `env:"API_POSTGRES_READ_TIMEOUT"  env-default:"10s"`
		WriteTimeout time.Duration `env:"API_POSTGRES_WRITE_TIMEOUT" env-default:"10s"`
		DisableTLS   bool          `env:"API_POSTGRES_DISABLE_TLS"   env-default:"true"`
	}
	Kaffka struct {
		Brokers  string `env:"KAFKA_BROKERS"       env-default:"0.0.0.0:9091"`
		Version  string `env:"KAFKA_VERSION"       env-default:"1.1.0"`
		Verbose  bool   `env:"KAFKA_VERBOSE"       env-default:"true"`
		ClientID string `env:"KAFKA_CLIENT_ID"     env-default:"sarama-easy"`
		Topics   string `env:"KAFKA_TOPICS"        env-default:""`

		TLSEnabled bool   `env:"KAFKA_TLS_ENABLED" env-default:"false"`
		TLSKey     string `env:"KAFKA_TLS_KEY"     env-default:""`
		TLSCert    string `env:"KAFKA_TLS_CERT"    env-default:""`
		CACerts    string `env:"KAFKA_CA_CERTS"    env-default:""`

		// Consumer specific parameters
		Group             string        `env:"KAFKA_GROUP"              env-default:"default-group"`
		RebalanceStrategy string        `env:"KAFKA_REBALANCE_STRATEGY" env-default:"roundrobin"`
		RebalanceTimeout  time.Duration `env:"KAFKA_REBALANCE_TIMEOUT"  env-default:"60s"`
		InitOffsets       string        `env:"KAFKA_INIT_OFFSETS"       env-default:"latest"`
		CommitInterval    time.Duration `env:"KAFKA_COMMIT_INTERVAL"    env-default:"10s"`

		// Producer specific parameters
		FlushInterval time.Duration `env:"KAFKA_FLUSH_INTERVAL"         env-default:"1s"`
	}
)

func (c Config) DatabaseStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&timezone=utc", c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Name)
}
