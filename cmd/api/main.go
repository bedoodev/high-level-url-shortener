// @title       URL Shortener API
// @version     1.0
// @description This is a simple high-level URL shortener written in Go.
// @host        localhost:8080
// @BasePath    /

package main

import (
	"context"
	"net/http"

	"github.com/bedoodev/high-level-url-shortener/internal/api"
	"github.com/bedoodev/high-level-url-shortener/internal/config"
	"github.com/bedoodev/high-level-url-shortener/internal/kafka"
	"github.com/bedoodev/high-level-url-shortener/internal/model"
	"github.com/bedoodev/high-level-url-shortener/internal/repository"
	"github.com/bedoodev/high-level-url-shortener/internal/service"
	"go.uber.org/zap"
)

func main() {
	if err := config.InitLogger(); err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}
	defer config.Logger.Sync()

	if err := config.InitPostgres(); err != nil {
		zap.L().Fatal("cannot initialize DB", zap.Error(err))
	}

	if err := config.DB.AutoMigrate(&model.URL{}); err != nil {
		zap.L().Fatal("failed to migrate", zap.Error(err))
	}

	config.InitCache()

	if err := config.InitRedis(); err != nil {
		zap.L().Fatal("failed to connect to Redis", zap.Error(err))
	}

	defer close(config.StopCleanupChan)

	repo := repository.NewURLRepository()
	ctx := context.Background()
	svc := service.NewURLService(repo)
	h := api.NewHandler(svc)
	r := api.NewRouter(h)

	kafka.InitKafkaProducer("kafka:9092")
	go kafka.StartClickConsumer(ctx, repo, "kafka:9092")

	zap.L().Info("Server is running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		zap.L().Fatal("server crashed", zap.Error(err))
	}
}
