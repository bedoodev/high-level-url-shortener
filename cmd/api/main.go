// @title       URL Shortener API
// @version     1.0
// @description This is a high-level URL shortener written in Go.
// @host        localhost:8080
// @BasePath    /

package main

import (
	"context"
	"net/http"
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
	// Initialize logger
	if err := config.InitLogger(); err != nil {
		panic("cannot initialize zap logger: " + err.Error())
	}
	defer config.Logger.Sync()

	// Setup context with cancel for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize PostgreSQL and run migrations
	if err := config.InitPostgres(); err != nil {
		zap.L().Fatal("cannot initialize DB", zap.Error(err))
	}
	if err := config.DB.AutoMigrate(&model.URL{}, &model.ClickEvent{}); err != nil {
		zap.L().Fatal("failed to migrate", zap.Error(err))
	}

	// Initialize local cache and Redis
	config.InitCache()
	if err := config.InitRedis(); err != nil {
		zap.L().Fatal("failed to connect to Redis", zap.Error(err))
	}
	defer close(config.StopCleanupChan)

	// Initialize Kafka
	brokerAddr := "kafka:9092"
	kafka.InitKafkaProducer(brokerAddr)

	repo := repository.NewURLRepository()

	// Start Kafka consumer in the background
	go kafka.StartClickConsumer(ctx, repo, brokerAddr)

	// Initialize HTTP handler and router
	svc := service.NewURLService(repo)
	h := api.NewHandler(svc)
	r := api.NewRouter(h)

	// Setup HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Run HTTP server in goroutine
	go func() {
		zap.L().Info("Server is running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("server crashed", zap.Error(err))
		}
	}()

	// Wait for SIGINT or SIGTERM
	<-ctx.Done()
	zap.L().Info("Shutdown signal received")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("failed to shutdown HTTP server", zap.Error(err))
	} else {
		zap.L().Info("HTTP server shutdown complete")
	}
}
