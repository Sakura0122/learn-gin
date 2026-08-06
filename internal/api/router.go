package api

import (
	"fmt"
	"net/http"

	"learn-gin/internal/api/article"
	"learn-gin/internal/api/user"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(ginMode string, appLogger *zap.Logger, db *gorm.DB) *gin.Engine {
	gin.SetMode(ginMode)

	// 创建一个名为 gin 的子 logger
	ginLogger := appLogger.Named("gin")
	// 接管 Gin 内部调试日志
	gin.DebugPrintFunc = func(format string, values ...any) {
		ginLogger.Debug(fmt.Sprintf(format, values...))
	}

	router := gin.New()
	// 异常恢复
	router.Use(ginzap.RecoveryWithZap(ginLogger, true))

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	api := router.Group("/api")

	user.RegisterRoutes(api, user.NewHandler(user.NewService(db)))
	article.RegisterRoutes(api, article.NewHandler(article.NewService(db)))

	return router
}
