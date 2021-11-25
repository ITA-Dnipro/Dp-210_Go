package main

import (
	"context"
	"database/sql"
	"fmt"
	cache "github.com/ITA-Dnipro/Dp-210_Go/auth/internal/cache/redis"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/client"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/usecase"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/repository/postgres"
	router "github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/http"

	grpcServer "github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/grpc"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/grpc/proto"

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

	db, err := sql.Open("pgx", cfg.DatabaseStr())
	if err != nil {
		log.Fatal(fmt.Errorf("creating db: %w", err))
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
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
	jwt, err := usecase.NewJwtAuth(cache.NewSessionCache(rdb, expire, cfg.TokenType), expire)
	if err != nil {
		log.Fatal(fmt.Errorf("initialize auth: %w", err))
	}

	kafka, err := client.NewKafka(cfg.KafkaBrokers, logger)
	if err != nil {
		log.Fatalf("init kafka: %w", err)
	}

	r, err := router.NewRouter(db, logger, kafka, rdb, jwt)
	_ = r
	if err != nil {
		log.Fatal(fmt.Errorf("initialize router: %w", err))
	}

	log.Println("Initialized successfully")

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.HttpPort), r); err != nil {
			log.Fatalf("failed to listen http: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen grpc: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterTokenValidatorServer(s, grpcServer.NewGrpcServer(jwt))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
