package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/handlers"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/kafkajob"
	repository "github.com/ITA-Dnipro/Dp-210_Go/payment/repo"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
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

	// DI
	repo := repository.NewRepository(db)
	handler := handlers.NewHandler(repo, logger)

	logger.Info("scheduling kafka produce")
	scheduler := gocron.NewScheduler(time.UTC)
	if _, err = scheduler.Every(3).Minute().Do(kafkajob.Produce, ctx, handler); err != nil {
		logger.Error("in func do of the scheduler error", zap.Error(err))
	}
	scheduler.StartAsync()

	logger.Info("kafka consume starting")
	go kafkajob.Consume(ctx, handler)

	logger.Fatal("server error", zap.Error(http.ListenAndServe(":1234", nil)))
}
