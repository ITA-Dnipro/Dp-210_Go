package config

import (
	"fmt"
	"net/url"
)

type Env struct {
	DbUser     string `env:"DB_USER" env-default:"postgres"`
	DbPassword string `env:"DB_PASSWORD" env-default:"secret"`
	DbName     string `env:"DB_NAME" env-default:"test"`
	DbHost     string `env:"DB_HOST" env-default:"0.0.0.0"`
	DbPort     string `env:"DB_PORT" env-default:"5432"`
	DbParams   string `env:"DB_PARAMS" env-default:"sslmode=disable&timezone=utc"`

	AppHost string `env:"APP_HOST" env-default:"localhost"`
	AppPort string `env:"APP_PORT" env-default:"8000"`
}

func (e *Env) ConnectionUrl() (*url.URL, error) {
	return url.Parse(e.ConnectionStr())
}

func (e *Env) ConnectionStr() string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?%v", e.DbUser, e.DbPassword, e.DbHost, e.DbPort, e.DbName, e.DbParams)

}
