package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/config"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres"
	router "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/server/http/middleware/auth"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/sender/mail"

	"github.com/go-redis/redis/v8"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

const (
	configPath     = "config.json"
	migrationsPath = "migrations"
)

// Main function
func main() {
	log.Println("Starting webapp dp210go")

	var env config.Env
	err := cleanenv.ReadEnv(&env)
	if err != nil {
		log.Fatal(fmt.Errorf("read env: %w", err))
	}

	var config config.Config
	err = cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatal(fmt.Errorf("read config: %w", err))
	}

	logger, _ := zap.NewProduction()

	gmail, err := mail.NewGmailEmailSender("config.json", "token.json")
	if err != nil {
		log.Fatal(fmt.Errorf("gmail sender: can't find files: %w", err))
	}

	jwtAuth, err := auth.NewAuthJwt()
	if err != nil {
		log.Fatal(fmt.Errorf("jwt auth: %w", err))
	}

	db, err := sql.Open("pgx", env.DatabaseStr())

	if err != nil {
		log.Fatal(fmt.Errorf("creating db: %w", err))
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal(fmt.Errorf("ping db %s : %w", env.DatabaseStr(), err))
	}

	err = postgres.MigrateUp(migrationsPath, env.DatabaseStr())
	if err != nil {
		log.Fatal(fmt.Errorf("db migrations: %w", err))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     env.RedisUrl,
		Password: env.RedisPassword,
		DB:       0,
	})

	_ = rdb

	r := router.NewRouter(db, logger, gmail, jwtAuth)
	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", env.AppHost, env.AppPort), r))
}
