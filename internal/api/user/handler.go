package user

import (
	"errors"
	"net/http"
	"strconv"

	"learn-gin/internal/common/result"

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
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, result.Error(result.CodeParamError, err.Error()))
		return
	}

	u, err := h.service.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUsernameExists):
			c.JSON(http.StatusOK, result.Error(result.CodeError, err.Error()))
		default:
			c.JSON(http.StatusOK, result.Error(result.CodeServerError, err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, result.Success(toUserResponse(u)))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, result.Error(result.CodeParamError, "invalid id"))
		return
	}

	u, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, result.Error(result.CodeNotFound, "user not found"))
			return
		}
		c.JSON(http.StatusOK, result.Error(result.CodeServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, result.Success(toUserResponse(u)))
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	list, total, err := h.service.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, result.Error(result.CodeServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":      toUserResponseList(list),
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, result.Error(result.CodeParamError, "invalid id"))
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, result.Error(result.CodeParamError, err.Error()))
		return
	}

	u, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, result.Error(result.CodeNotFound, "user not found"))
			return
		}
		c.JSON(http.StatusOK, result.Error(result.CodeServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, result.Success(toUserResponse(u)))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, result.Error(result.CodeParamError, "invalid id"))
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, result.Error(result.CodeNotFound, "user not found"))
			return
		}
		c.JSON(http.StatusOK, result.Error(result.CodeServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, result.Success(nil))
}
