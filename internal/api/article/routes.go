package article

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	r.POST("/articles", handler.Create)
	r.GET("/articles/:id", handler.GetByID)
	r.GET("/articles", handler.List)
	r.PUT("/articles/:id", handler.Update)
	r.DELETE("/articles/:id", handler.Delete)
}
