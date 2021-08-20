package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ITA-Dnipro/Dp-210_Go/config"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres"
	router "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http"

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

	r := router.NewRouter(db, logger)
	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", env.AppHost, env.AppPort), r))
}
