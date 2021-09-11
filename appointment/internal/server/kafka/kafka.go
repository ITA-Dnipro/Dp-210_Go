package kafka

import (
	"encoding/json"

	"github.com/ITA-Dnipro/Dp-210_Go/appointment/internal/entity"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

var (
	AppoinmentTopic   = "appointment"
	NotificationTopic = "notification"
	BillTopic         = "bill"
)

type Kafka struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
	logger   *zap.Logger
}

type fnHandler func(payload []byte) error

func NewKafka(brokers []string, logger *zap.Logger) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Kafka{Producer: producer, Consumer: consumer, logger: logger}, nil
}

func (k *Kafka) send(topic string, payload interface{}) error {
	msg, err := json.Marshal(payload)
	if err != nil {
		k.logger.Error("can't marsha message", zap.String("topic", topic), zap.Error(err))
		return err
	}
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
	p, o, err := k.Producer.SendMessage(message)
	k.logger.Info(
		"send message",
		zap.String("topic", topic),
		zap.Int32("partition", p),
		zap.Int64("offser", o),
	)
	return err
}

func (k *Kafka) SendAppointment(a *entity.Appointment) error {
	return k.send(AppoinmentTopic, a)
}
func (k *Kafka) SendBill(a *entity.Appointment) error {
	return k.send(AppoinmentTopic, a)
}
func (k *Kafka) SendNotification(a *entity.Appointment) error {
	return k.send(AppoinmentTopic, a)
}

func (k *Kafka) on(topic string, handler fnHandler) error {
	partitions, err := k.Consumer.Partitions(topic)
	if err != nil {
		return err
	}
	for _, partition := range partitions {
		pc, err := k.Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func() {
			for {
				message := <-pc.Messages()
				k.logger.Info("recive message ", zap.String("topic", message.Topic))
				if err := handler(message.Value); err != nil {
					k.logger.Error("handel message",
						zap.String("topic", message.Topic),
						zap.Error(err),
					)
				}
			}
		}()
	}
	return nil
}

func (k *Kafka) OnAppointment(handler fnHandler) error {
	return k.on(AppoinmentTopic, handler)
}

func (k *Kafka) Close() {
	k.Consumer.Close()
	k.Producer.Close()
}
