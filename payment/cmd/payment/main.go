package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/handlers/kafkahand"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/kafka"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/proto/statistics"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/repository"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"

	// This is a database/sql compatibility layer for pgx. pgx can be used as a normal database/sql
	//driver, but at any time, the native interface can be acquired for more performance or PostgreSQL
	//specific functionality.
	_ "github.com/jackc/pgx/v4/stdlib"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

const (
	user     = "postgres"
	password = "0310"
	host     = "localhost"
	port     = "5432"
	dbname   = "paymentdb"
	params   = "sslmode=disable&timezone=utc"
)

var strConn = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?%v", user, password, host, port, dbname, params)

func main() {
	logger, _ := zap.NewProduction()
	ctx := context.Background()

	logger.Info("db preparing")
	db, err := sql.Open("pgx", strConn)
	if err != nil {
		logger.Fatal("creating db:", zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		if err = db.Close(); err != nil {
			logger.Error("db close error", zap.Error(err))
		}
		logger.Fatal("db ping failed", zap.Error(err))
	}

	conn, err := grpc.Dial(":1235", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	statClient := statistics.NewStatClient(conn)

	// DI
	repo := repository.NewRepository(db)
	kafkaHandler := kafkahand.NewHandler(repo, logger, statClient)

	logger.Info("scheduling kafka produce")
	scheduler := gocron.NewScheduler(time.UTC)
	if _, err = scheduler.Every(1).Minute().Do(kafka.Produce, ctx, kafkaHandler); err != nil {
		logger.Error("in func do of the scheduler error", zap.Error(err))
	}
	scheduler.StartAsync()

	logger.Info("kafka consume starting")
	go kafka.Consume(ctx, kafkaHandler)

	logger.Fatal("server error", zap.Error(http.ListenAndServe(":1234", nil)))
}
