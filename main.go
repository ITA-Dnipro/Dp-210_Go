package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ITA-Dnipro/Dp-210_Go/config"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/repository/postgres"
	router "github.com/ITA-Dnipro/Dp-210_Go/internal/server/http"
	"github.com/ITA-Dnipro/Dp-210_Go/internal/test"

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

	test.InitTestData(db)

	logger.Info("starting web server")
	r := router.NewRouter(db, logger)
	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", env.AppHost, env.AppPort), r))
}
