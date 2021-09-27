package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

const (
	Topic = "notification"
)

type fnHandler func(payload []byte) error

type Kafka struct {
	Consumer sarama.Consumer
	logger   *zap.Logger
}

func NewKafka(brokers []string, logger *zap.Logger) (*Kafka, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("init kafka: %w", err)
	}

	return &Kafka{consumer, logger}, nil
}

func (k *Kafka) on(topic string, handler fnHandler) error {
	partitions, err := k.Consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("kafka on: %w", err)
	}

	for _, partition := range partitions {
		pc, err := k.Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return fmt.Errorf("kafka consume: %w", err)
		}

		go func() {
			for {
				message := <-pc.Messages()
				t := zap.String("topic", topic)
				k.logger.Info("received message ", t)
				if err := handler(message.Value); err != nil {
					k.logger.Error("handle message ", t, zap.Error(err))
				}
			}
		}()
	}

	return nil
}

func (k *Kafka) OnEmail(handler fnHandler) error {
	return k.on(Topic, handler)
}

func (k *Kafka) Close() {
	k.Consumer.Close()
}
