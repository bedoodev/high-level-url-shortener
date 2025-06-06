package kafka

import (
	"context"

	"github.com/bedoodev/high-level-url-shortener/internal/repository"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func StartClickConsumer(ctx context.Context, repo repository.URLRepository, brokerAddr string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddr},
		Topic:   "click-events",
		GroupID: "click-worker-group",
	})

	zap.L().Info("Kafka consumer started...")

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			zap.L().Error("Error reading Kafka message", zap.Error(err))
			continue
		}

		shortCode := string(m.Value)
		zap.L().Info("Received click event", zap.String("shortCode", shortCode))

		if err := repo.IncrementClickCount(ctx, shortCode); err != nil {
			zap.L().Error("Failed to update click count", zap.Error(err))
		} else {
			zap.L().Info("Click count incremented", zap.String("shortCode", shortCode))
		}
	}
}
