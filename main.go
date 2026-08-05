package main

import (
	"log"
	"net/http"

	"learn-gin/internal/api/article"
	"learn-gin/internal/api/user"
	"learn-gin/internal/config"
	"learn-gin/internal/infra/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}
	log.Printf("database connected: %s", cfg.DSN())

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := router.Group("/api/v1")
	user.RegisterRoutes(api, user.NewHandler(user.NewService(db)))
	article.RegisterRoutes(api, article.NewHandler(article.NewService(db)))

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
