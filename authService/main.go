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
	log.Println("Starting webapp dp210go")

	var env config.Config
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

	gmail, err := sender.NewGmailEmailSender("config.json", "token.json")
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

	r, err := router.NewRouter(db, logger, gmail, rdb)
	if err != nil {
		log.Fatal(fmt.Errorf("initialize router: %w", err))
	}
	// Start server
	log.Println("Initialized successfully")

	if err := http.ListenAndServe(fmt.Sprintf(":%v", env.AppPort), r); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", string(8081))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	expire := time.Minute * 15
	jwt, err := usecase.NewJwtAuth(cache.NewSessionCache(rdb, expire, "jwtToken"), expire)

	proto.RegisterTokenValidatorServer(s, grpcServer.NewGrpcServer(jwt))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
