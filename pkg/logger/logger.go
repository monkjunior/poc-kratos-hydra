package logger

import (
	"github.com/monkjunior/poc-kratos-hydra/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

func getLogConfig() zap.Config {
	zapCfg := zap.NewDevelopmentConfig()

	logMode := config.Cfg.Log.Mode
	if logMode == "prod" || logMode == "production" {
		zapCfg = zap.NewProductionConfig()
	}

	logLevel := config.Cfg.Log.Level
	switch logLevel {
	case "debug":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "d-panic":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.DPanicLevel)
	case "panic":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	case "fatal":
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		zapCfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	return zapCfg
}

func InitLogger() {
	cfg := getLogConfig()
	Logger, _ = cfg.Build()
}

func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}
