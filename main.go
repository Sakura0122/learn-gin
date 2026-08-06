package main

import (
	"fmt"

	"learn-gin/internal/api"
	"learn-gin/internal/config"
	"learn-gin/internal/infra/database"
	"learn-gin/internal/infra/logging"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	appLogger, err := logging.New(cfg.AppEnv, cfg.LogDir)
	if err != nil {
		panic(fmt.Sprintf("初始化日志失败：%v", err))
	}
	defer appLogger.Sync()

	db, err := database.New(cfg, appLogger)
	if err != nil {
		appLogger.Fatal("连接数据库失败", zap.Error(err))
	}
	appLogger.Info("数据库连接成功")

	router := api.NewRouter(cfg.GinMode, appLogger, db)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		appLogger.Fatal("服务启动失败", zap.Error(err))
	}
}
