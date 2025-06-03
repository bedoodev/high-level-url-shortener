package config

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger() error {
	var err error

	// Production logger
	Logger, err = zap.NewProduction()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(Logger)
	return nil
}
