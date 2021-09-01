package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/config"
	cache "github.com/ITA-Dnipro/Dp-210_Go/internal/cache/redis"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres"
	router "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/service/auth"
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

	var cfg config.Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("read config: %w", err))
	}

	logger, _ := zap.NewProduction()

	gmail, err := mail.NewGmailEmailSender("config.json", "token.json")
	if err != nil {
		log.Fatal(fmt.Errorf("gmail sender: can't find files: %w", err))
	}

	db, err := sql.Open("pgx", env.DatabaseStr())
	if err != nil {
		log.Fatal(fmt.Errorf("creating db: %w", err))
	}

	err = db.Ping()
	if err != nil {
		if err = db.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		}
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

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(fmt.Errorf("connect to redis server: %w", err))
	}

	tokenExp := 15 * time.Minute
	jwtAuth, err := auth.NewJwtAuth(cache.NewSessionCache(rdb, tokenExp, "jwtToken"), tokenExp)
	if err != nil {
		log.Fatal(fmt.Errorf("jwt auth: %w", err))
	}

	r := router.NewRouter(db, logger, gmail, jwtAuth, rdb)
	// Start server
	log.Println("Initialized successfully")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", env.AppPort), r))
}
