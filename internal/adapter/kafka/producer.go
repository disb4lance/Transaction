package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"
)

type TransactionCreatedEvent struct {
	TransactionID string          `json:"transaction_id"`
	UserID        string          `json:"user_id"`
	CategoryID    string          `json:"category_id"`
	Amount        decimal.Decimal `json:"amount"`
	CreatedAt     time.Time       `json:"created_at"`
}

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

func (p *Producer) PublishTransactionCreated(
	ctx context.Context,
	event TransactionCreatedEvent,
) error {

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.UserID),
		Value: payload,
	})
}
