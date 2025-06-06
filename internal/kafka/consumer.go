package kafka

import (
	"context"

	"github.com/bedoodev/high-level-url-shortener/internal/config"
	"github.com/bedoodev/high-level-url-shortener/internal/model"
	"github.com/bedoodev/high-level-url-shortener/internal/repository"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func StartClickConsumer(ctx context.Context, repo repository.URLRepository, brokerAddr string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddr},
		Topic:    "click-events",
		GroupID:  "url-shortener-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	go func() {
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				zap.L().Error("Failed to read Kafka message", zap.Error(err))
				continue
			}

			shortCode := string(m.Value)
			zap.L().Info("Received click event", zap.String("shortCode", shortCode))

			// 🔍 shortCode ile URL'yi bul
			url, err := repo.FindByShortCode(ctx, shortCode)
			if err != nil {
				zap.L().Error("Failed to find URL", zap.String("shortCode", shortCode), zap.Error(err))
				continue
			}

			// 📝 click_events tablosuna yaz
			click := &model.ClickEvent{
				URLID: url.ID,
			}
			if err := config.DB.Create(click).Error; err != nil {
				zap.L().Error("Failed to insert click event", zap.Error(err))
			}
		}
	}()
}
