package config

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() error {
	var err error

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	// Production logger
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(Logger)
	return nil
}
