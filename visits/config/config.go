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
		Brokers  string `envconfig:"KAFKA_BROKERS"       env-default:"0.0.0.0:9092"`
		Version  string `envconfig:"KAFKA_VERSION"       env-default:"1.1.0"`
		Verbose  bool   `envconfig:"KAFKA_VERBOSE"       env-default:"true"`
		ClientID string `envconfig:"KAFKA_CLIENT_ID"     env-default:"sarama-easy"`
		Topics   string `envconfig:"KAFKA_TOPICS"        env-default:""`

		TLSEnabled bool   `envconfig:"KAFKA_TLS_ENABLED" env-default:"false"`
		TLSKey     string `envconfig:"KAFKA_TLS_KEY"     env-default:""`
		TLSCert    string `envconfig:"KAFKA_TLS_CERT"    env-default:""`
		CACerts    string `envconfig:"KAFKA_CA_CERTS"    env-default:""`

		// Consumer specific parameters
		Group             string        `envconfig:"KAFKA_GROUP"              env-default:"default-group"`
		RebalanceStrategy string        `envconfig:"KAFKA_REBALANCE_STRATEGY" env-default:"roundrobin"`
		RebalanceTimeout  time.Duration `envconfig:"KAFKA_REBALANCE_TIMEOUT"  env-default:"60s"`
		InitOffsets       string        `envconfig:"KAFKA_INIT_OFFSETS"       env-default:"latest"`
		CommitInterval    time.Duration `envconfig:"KAFKA_COMMIT_INTERVAL"    env-default:"10s"`

		// Producer specific parameters
		FlushInterval time.Duration `envconfig:"KAFKA_FLUSH_INTERVAL"         env-default:"1s"`
	}
)
