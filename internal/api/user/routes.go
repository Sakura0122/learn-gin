package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	users := r.Group("/users")

	users.POST("", handler.Create)
	users.GET("/:id", handler.GetByID)
	users.GET("", handler.List)
	users.PUT("/:id", handler.Update)
	users.DELETE("/:id", handler.Delete)
}
