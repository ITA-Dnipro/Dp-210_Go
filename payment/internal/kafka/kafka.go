package kafka

import (
	"context"
	"fmt"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/entity"
	"github.com/ITA-Dnipro/Dp-210_Go/payment/internal/handlers/kafkahand"

	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	topic         = "bill"
	topicOut      = "bill-period"
	broker        = "localhost:9092"
	consumerGroup = "bill-group"
)

var w = kafkago.Writer{
	Addr:  kafkago.TCP(broker),
	Topic: topicOut,
}

func Produce(ctx context.Context, h *kafkahand.Handler) {
	report, err := h.SendMonthlyReport()
	if err != nil {
		h.Logger.Error("send monthly report to topic info >", zap.Error(err))
		return
	}

	if report == nil {
		return
	}

	if err = w.WriteMessages(ctx, kafkago.Message{
		Key:   []byte("to-bill-period"),
		Value: report,
	}); err != nil {
		h.Logger.Error("could not write message error", zap.Error(err))
		return
	}

	h.Logger.Info(fmt.Sprintf("\n>>> SENT TO KAFKA %v\n", string(report)))
}

func Consume(ctx context.Context, h *kafkahand.Handler) {
	var r = kafkago.NewReader(kafkago.ReaderConfig{
		Topic:   topic,
		Brokers: []string{broker},
		GroupID: consumerGroup,
	})
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			h.Logger.Error("could not read message", zap.Error(err))
			continue
		}

		bill, err := entity.NewBill(msg.Value)
		if err != nil {
			h.Logger.Error("new data error", zap.Error(err))
			continue
		}

		if err = h.InsertToDb(bill); err != nil {
			h.Logger.Error("insert to db error", zap.Error(err))
		}
	}
}
