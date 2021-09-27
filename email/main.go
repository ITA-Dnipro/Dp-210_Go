package main

import (
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/email/internal/config"
	"github.com/ITA-Dnipro/Dp-210_Go/email/internal/server/kafka"
	"github.com/ITA-Dnipro/Dp-210_Go/email/internal/usesase"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"log"
	"os"
)

const configPath = "config.json"

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("building logger", err)
	}

	var cfg config.Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("read config: %w", err))
	}

	logger.Info("Initializing kafka")
	k, err := kafka.NewKafka(cfg.KafkaBrokers, logger)

	if err != nil {
		logger.Error("error:", zap.Error(fmt.Errorf("connecting to kafka: %w", err)))
		os.Exit(1)
	}
	defer k.Close()

	mail, err := usesase.NewGmailEmailSender("config.json", "token.json")
	if err != nil {
		log.Fatal(err)
	}

	h := kafka.EventHandler{Sender: mail}

	k.OnEmail(h.EmailFromEvent)
}
