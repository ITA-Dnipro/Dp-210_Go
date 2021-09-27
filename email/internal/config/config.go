package config

type Config struct {
	KafkaBrokers []string `env:"KAFKA_BROKERS" env-default:"0.0.0.0:9091"`
}
