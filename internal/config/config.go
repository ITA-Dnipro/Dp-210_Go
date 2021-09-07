package config

import (
	"fmt"
	"net/url"
	"sync"
)

var conf Config
var mu sync.Mutex

type Config struct {
	DbUser     string `env:"DB_USER" env-default:"postgres"`
	DbPassword string `env:"POSTGRES_PASSWORD" env-default:"dp210go"`
	DbName     string `env:"DB_NAME" env-default:"postgres"`
	DbHost     string `env:"DB_HOST" env-default:"db"`
	DbPort     string `env:"DB_PORT" env-default:"5432"`
	DbParams   string `env:"DB_PARAMS" env-default:"sslmode=disable&timezone=utc"`

	RedisUrl      string `env:"REDIS_URL" env-default:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" env-default:""`

	AppHost string `env:"APP_HOST" env-default:"localhost"`
	AppPort string `env:"APP_PORT" env-default:"8000"`

	AuthAddress string `json:"auth_service_address" end-default:"localhost:8002"`
}

func (e *Config) DatabaseUrl() (*url.URL, error) {
	return url.Parse(e.DatabaseStr())
}

func (e *Config) DatabaseStr() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?%v", e.DbUser, e.DbPassword, e.DbHost, e.DbPort, e.DbName, e.DbParams)

}

func GetConfig() Config {
	return conf
}

func SetConfig(cfg Config) {
	mu.Lock()
	conf = cfg
	mu.Unlock()
}
