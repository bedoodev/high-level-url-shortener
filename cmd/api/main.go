// @title       URL Shortener API
// @version     1.0
// @description This is a high-level URL shortener written in Go.
// @host        localhost:8080
// @BasePath    /

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bedoodev/high-level-url-shortener/internal/api"
	"github.com/bedoodev/high-level-url-shortener/internal/config"
	"github.com/bedoodev/high-level-url-shortener/internal/kafka"
	"github.com/bedoodev/high-level-url-shortener/internal/model"
	"github.com/bedoodev/high-level-url-shortener/internal/repository"
	"github.com/bedoodev/high-level-url-shortener/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Logger
	if err := config.InitLogger(); err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}
	defer config.Logger.Sync()

	// Signal & Context Setup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown için sinyal yakalayıcı
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Database
	if err := config.InitPostgres(); err != nil {
		zap.L().Fatal("cannot initialize DB", zap.Error(err))
	}
	if err := config.DB.AutoMigrate(&model.URL{}); err != nil {
		zap.L().Fatal("failed to migrate", zap.Error(err))
	}

	// Caches
	config.InitCache()
	if err := config.InitRedis(); err != nil {
		zap.L().Fatal("failed to connect to Redis", zap.Error(err))
	}
	defer close(config.StopCleanupChan)

	// Kafka
	brokerAddr := "kafka:9092"
	kafka.InitKafkaProducer(brokerAddr)

	repo := repository.NewURLRepository()

	// Kafka consumer'ı arka planda çalıştır
	go kafka.StartClickConsumer(ctx, repo, brokerAddr)

	// API servis
	svc := service.NewURLService(repo)
	h := api.NewHandler(svc)
	r := api.NewRouter(h)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// HTTP server'ı arka planda başlat
	go func() {
		zap.L().Info("Server is running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("server crashed", zap.Error(err))
		}
	}()

	// Sinyal gelene kadar bekle
	<-signalChan
	zap.L().Info("Shutdown signal received")

	// Graceful shutdown başlat
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("failed to shutdown HTTP server", zap.Error(err))
	} else {
		zap.L().Info("HTTP server shutdown complete")
	}

	cancel() // Kafka gibi işlemleri iptal et
}
