package api

import (
	"learn-gin/internal/api/article"
	"learn-gin/internal/api/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")

	user.RegisterRoutes(api, user.NewHandler(user.NewService(db)))
	article.RegisterRoutes(api, article.NewHandler(article.NewService(db)))
}
