package client

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

const (
	PasswCodeTopic = "password_code"
)

type fnHandler func(payload []byte) error

type Kafka struct {
	Producer sarama.SyncProducer
	logger   *zap.Logger
}

func NewKafka(brokers []string, logger *zap.Logger) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Kafka{Producer: producer, logger: logger}, nil
}

func (k *Kafka) Send(topic string, payload interface{}) error {
	msg, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("send to kafka: marshal: %w", err)
	}

	m := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	part, off, err := k.Producer.SendMessage(m)
	k.logger.Info(
		"Send to kafka",
		zap.String("topic: ", topic),
		zap.Int32("partition: ", part),
		zap.Int64("offset: ", off),
	)

	return nil
}
