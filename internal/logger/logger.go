package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
	var cfg zap.Config

	if os.Getenv("APP_ENV") == "development" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	return zapLogger
}
