package main

import (
	"context"
	"database/sql"
	"fmt"
	cache "github.com/ITA-Dnipro/Dp-210_Go/authService/internal/cache/redis"
	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/usecase"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/repository/postgres"
	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/sender"
	router "github.com/ITA-Dnipro/Dp-210_Go/authService/internal/server/http"

	grpcServer "github.com/ITA-Dnipro/Dp-210_Go/authService/internal/server/grpc"
	"github.com/ITA-Dnipro/Dp-210_Go/authService/internal/server/grpc/proto"

	"github.com/go-redis/redis/v8"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

const (
	configPath     = "config.json"
	migrationsPath = "migrations"
)

func main() {
	log.Println("Starting authentication microservice")

	var cfg config.Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("read config: %w", err))
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("read env: %w", err))
	}

	config.SetConfig(cfg)

	logger, _ := zap.NewProduction()

	gmail, err := sender.NewGmailEmailSender(configPath, "token.json")
	if err != nil {
		log.Fatal(fmt.Errorf("gmail sender: can't find files: %w", err))
	}

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

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(fmt.Errorf("connect to redis server: %w", err))
	}

	expire := time.Duration(cfg.TokenExpirationMillis) * time.Millisecond
	jwt, err := usecase.NewJwtAuth(cache.NewSessionCache(rdb, expire, "jwtToken"), expire)
	if err != nil {
		log.Fatal(fmt.Errorf("initialize auth: %w", err))
	}

	r, err := router.NewRouter(db, logger, gmail, rdb, jwt)
	if err != nil {
		log.Fatal(fmt.Errorf("initialize router: %w", err))
	}

	log.Println("Initialized successfully")

	if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.HttpPort), r); err != nil {
		log.Fatalf("failed to listen http: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen grpc: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterTokenValidatorServer(s, grpcServer.NewGrpcServer(jwt))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
