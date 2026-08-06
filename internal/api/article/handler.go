package article

import (
	"errors"
	"learn-gin/internal/common"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.Error(common.CodeParamError, err.Error()))
		return
	}

	a, err := h.service.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			c.JSON(http.StatusOK, common.Error(common.CodeNotFound, err.Error()))
		default:
			c.JSON(http.StatusOK, common.Error(common.CodeServerError, err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, common.Success(toArticleResponse(a)))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, common.Error(common.CodeParamError, "invalid id"))
		return
	}

	a, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.Error(common.CodeNotFound, "article not found"))
			return
		}
		c.JSON(http.StatusOK, common.Error(common.CodeServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(toArticleResponse(a)))
}

func (h *Handler) List(c *gin.Context) {
	userID, _ := uuid.Parse(c.Query("user_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	list, total, err := h.service.List(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, common.Error(common.CodeServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(gin.H{
		"list":      toArticleResponseList(list),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, common.Error(common.CodeParamError, "invalid id"))
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, common.Error(common.CodeParamError, err.Error()))
		return
	}

	a, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.Error(common.CodeNotFound, "article not found"))
			return
		}
		c.JSON(http.StatusOK, common.Error(common.CodeServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(toArticleResponse(a)))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, common.Error(common.CodeParamError, "invalid id"))
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.Error(common.CodeNotFound, "article not found"))
			return
		}
		c.JSON(http.StatusOK, common.Error(common.CodeServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, common.Success(nil))
}
