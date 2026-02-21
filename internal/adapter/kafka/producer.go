package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"transaction-service/internal/adapter/kafka/kmodel"
	"transaction-service/internal/config"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducerWithBrokers(cfg config.KafkaConfig) *Producer {
	brokers := strings.Split(cfg.Brokers, ",")

	for i, b := range brokers {
		brokers[i] = strings.TrimSpace(b)
	}

	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
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
