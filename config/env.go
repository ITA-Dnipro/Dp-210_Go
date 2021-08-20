package config

import (
	"fmt"
	"net/url"
)

type Env struct {
	DbUser     string `env:"DB_USER" env-default:"postgres"`
	DbPassword string `env:"POSTGRES_PASSWORD" env-default:"dp210go"`
	DbName     string `env:"DB_NAME" env-default:"postgres"`
	DbHost     string `env:"DB_HOST" env-default:"localhost"`
	DbPort     string `env:"DB_PORT" env-default:"5432"`
	DbParams   string `env:"DB_PARAMS" env-default:"sslmode=disable&timezone=utc"`

	RedisUrl      string `env:"REDIS_URL" env-default:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" env-default:""`

	AppHost string `env:"APP_HOST" env-default:"localhost"`
	AppPort string `env:"APP_PORT" env-default:"8000"`
}

func (e *Env) DatabaseUrl() (*url.URL, error) {
	return url.Parse(e.DatabaseStr())
}

func (e *Env) DatabaseStr() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?%v", e.DbUser, e.DbPassword, e.DbHost, e.DbPort, e.DbName, e.DbParams)

}
