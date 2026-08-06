package database

import (
	"learn-gin/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func New(cfg *config.Config, appLogger *zap.Logger) (*gorm.DB, error) {
	logLevel := logger.Info
	if cfg.AppEnv == "production" {
		logLevel = logger.Warn
	}

	// 创建 GORM 专用的 Zap 适配器
	dbLogger := zapgorm2.New(appLogger.Named("gorm"))
	dbLogger.LogLevel = logLevel

	return gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: dbLogger,
	})
}
