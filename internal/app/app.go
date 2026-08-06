package app

import (
	"fmt"
	"net/http"

	"learn-gin/internal/api"
	"learn-gin/internal/config"
	"learn-gin/internal/infra/database"
	"learn-gin/internal/infra/logging"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application struct {
	config *config.Config
	logger *zap.Logger
	router *gin.Engine
}

func New() *Application {
	cfg := config.Load()

	appLogger, err := logging.New(cfg.AppEnv, cfg.LogDir)
	if err != nil {
		panic(fmt.Sprintf("初始化日志失败：%v", err))
	}

	db, err := database.New(cfg, appLogger)
	if err != nil {
		appLogger.Fatal("连接数据库失败", zap.Error(err))
	}
	appLogger.Info("数据库连接成功")

	return &Application{
		config: cfg,
		logger: appLogger,
		router: newRouter(cfg.GinMode, appLogger, db),
	}
}

func (app *Application) Run() {
	defer app.logger.Sync()

	app.router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	app.logger.Info(
		"服务开始启动",
		zap.String("url", "http://localhost:"+app.config.ServerPort),
	)

	err := app.router.Run(":" + app.config.ServerPort)
	if err != nil {
		app.logger.Fatal("服务启动失败", zap.Error(err))
	}
}

func newRouter(ginMode string, appLogger *zap.Logger, db *gorm.DB) *gin.Engine {
	gin.SetMode(ginMode)

	ginLogger := appLogger.Named("gin")
	gin.DebugPrintFunc = func(format string, values ...any) {
		ginLogger.Debug(fmt.Sprintf(format, values...))
	}

	router := gin.New()
	router.Use(ginzap.RecoveryWithZap(ginLogger, true))
	api.RegisterRoutes(router, db)

	return router
}
