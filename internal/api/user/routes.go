package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	r.POST("/users", handler.Create)
	r.GET("/users/:id", handler.GetByID)
	r.GET("/users", handler.List)
	r.PUT("/users/:id", handler.Update)
	r.DELETE("/users/:id", handler.Delete)
}
