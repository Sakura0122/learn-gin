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
		log.Fatalf("connect database failed: %v", err)
	}
	log.Printf("database connected: %s", cfg.DSN())

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api.RegisterRoutes(router, db)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
