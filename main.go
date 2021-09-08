package main

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"

	"github.com/ITA-Dnipro/Dp-210_Go/internal/config"
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

func main() {
	log.Println("Starting webapp dp210go")

	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("read env: %w", err))
	}

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("read config: %w", err))
	}

	config.SetConfig(cfg)

	logger, _ := zap.NewProduction()

	db, err := sql.Open("pgx", cfg.DatabaseStr())
	if err != nil {
		log.Fatal(fmt.Errorf("creating db: %w", err))
	}

	err = db.Ping()
	if err != nil {
		if err = db.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		log.Fatal(fmt.Errorf("ping db %s : %w", cfg.DatabaseStr(), err))
	}

	err = postgres.MigrateUp(migrationsPath, cfg.DatabaseStr())
	if err != nil {
		log.Fatal(fmt.Errorf("db migrations: %w", err))
	}

	//move to main, inject connection
	conn, err := grpc.Dial(config.GetConfig().AuthAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not validate token via grpc %v", err)
		return
	}
	defer conn.Close()

	r := router.NewRouter(db, logger, conn)
	// Start server
	log.Println("Initialized successfully")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", cfg.AppPort), r))
}
