package main

import (
	"net/http"

	"github.com/bedoodev/high-level-url-shortener/internal/config"
	"github.com/bedoodev/high-level-url-shortener/internal/model"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
		zap.L().Info("request received", zap.String("method", r.Method), zap.String("path", r.URL.Path))
	})

	zap.L().Info("Server is starting", zap.String("port", ":8080"))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		zap.L().Fatal("server failed", zap.Error(err))
	}
}
