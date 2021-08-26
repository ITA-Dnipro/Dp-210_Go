package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
)

var AppoinmentTopic = "appointment"
var MailTopic = "mail"
var BillTopic = "bill"

type Events interface {
	Emit(topic string, payload interface{}) error
	On(topic string, handler fnHandler) error
}

type Kafka struct {
	Emitter  sarama.SyncProducer
	Listener sarama.Consumer
}

type fnHandler func(payload []byte) error

func NewEvents(brokers []string) (Events, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	//config.Producer.Partitioner = sarama.NewRandomPartitioner
	//config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// config.Consumer.Return.Errors = true
	// config.Consumer.Offsets.CommitInterval = 10 * time.Second
	// config.Consumer.Group.Rebalance.Timeout = 1 * time.Minute
	// config.Consumer.Group.Rebalance.Retry.Max = 6
	// config.Consumer.Group.Rebalance.Retry.Backoff = 2 * time.Second
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	// config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Kafka{producer, consumer}, nil
}

func (k *Kafka) Emit(topic string, payload interface{}) error {
	msg, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	_, _, err = k.Emitter.SendMessage(message)
	return err
}

func (k *Kafka) On(topic string, handler fnHandler) error {
	partitions, err := k.Listener.Partitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, err := k.Listener.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func() {
			for {
				message := <-pc.Messages()
				handler(message.Value)
			}
		}()
	}

	return nil
}
