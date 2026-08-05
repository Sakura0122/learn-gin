package article

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	articles := r.Group("/articles")

	articles.POST("", handler.Create)
	articles.GET("/:id", handler.GetByID)
	articles.GET("", handler.List)
	articles.PUT("/:id", handler.Update)
	articles.DELETE("/:id", handler.Delete)
}
