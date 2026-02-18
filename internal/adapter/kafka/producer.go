package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"transaction-service/internal/adapter/kafka/kmodel"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(
	broker string,
	topic string,
) *Producer {

	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishTransactionEvent(
	ctx context.Context,
	topic string,
	event kmodel.TransactionEvent,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(event.UserID),
		Value: payload,
	})
}
