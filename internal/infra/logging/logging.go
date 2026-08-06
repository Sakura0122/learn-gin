package logging

import (
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(appEnv, logDir string) (*zap.Logger, error) {
	var config zap.Config
	if appEnv == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleLogger, err := config.Build()
	if err != nil {
		return nil, err
	}
	fileWriter, err := rotatelogs.New(
		filepath.Join(logDir, "app-%Y-%m-%d.log"),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithRotationCount(7),
	)
	if err != nil {
		return nil, err
	}
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(config.EncoderConfig),
		zapcore.AddSync(fileWriter),
		config.Level,
	)

	return consoleLogger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, fileCore)
	})), nil
}
