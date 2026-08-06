package main

import (
	"log"
	"net/http"

	"learn-gin/internal/api"
	"learn-gin/internal/config"
	"learn-gin/internal/infra/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	log.Print("数据库连接成功")

	router := gin.Default()
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	api.RegisterRoutes(router, db)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("服务启动失败：%v", err)
	}
}
