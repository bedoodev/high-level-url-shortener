package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var writer *kafka.Writer

func InitKafkaProducer(brokerAddr string) {
	writer = kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:  []string{brokerAddr},
			Topic:    "click-events",
			Balancer: &kafka.LeastBytes{},
		})

	zap.L().Info("Kafka producer initialized", zap.String("broker", brokerAddr))
}

// PublishClick sends shortCode to Kafka
func PublishClick(ctx context.Context, shortCode string) error {
	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(shortCode),
		Value: []byte(shortCode),
		Time:  time.Now(),
	})

	if err != nil {
		zap.L().Error("Failed to publish click event", zap.Error(err))
	}

	zap.L().Info("Published click event", zap.String("shortCode", shortCode)) // Log the shortCode for tracing purpo

	return err
}
