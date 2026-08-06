package article

import (
	"errors"
	"learn-gin/internal/common/page"
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
	var req CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(result.CodeParamError, "请求参数错误", c)
		return
	}

	a, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound):
			result.Error(result.CodeNotFound, "用户不存在", c)
		default:
			result.Error(result.CodeServerError, "服务器内部错误", c)
		}
		return
	}
	result.Success(toArticleResponse(a), c)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		result.Error(result.CodeParamError, "ID格式错误", c)
		return
	}

	a, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Error(result.CodeNotFound, "文章不存在", c)
			return
		}
		result.Error(result.CodeServerError, "服务器内部错误", c)
		return
	}
	result.Success(toArticleResponse(a), c)
}

func (h *Handler) List(c *gin.Context) {
	userID, _ := uuid.Parse(c.Query("user_id"))
	pageRequest, err := page.Parse(c)
	if err != nil {
		result.Error(result.CodeParamError, "分页参数错误", c)
		return
	}

	list, total, err := h.service.List(c.Request.Context(), userID, pageRequest)
	if err != nil {
		if errors.Is(err, page.ErrUnsupportedSortField) {
			result.Error(result.CodeParamError, err.Error(), c)
			return
		}
		result.Error(result.CodeServerError, "服务器内部错误", c)
		return
	}
	result.Success(page.NewResult(pageRequest, total, toArticleResponseList(list)), c)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		result.Error(result.CodeParamError, "ID格式错误", c)
		return
	}

	var req UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(result.CodeParamError, "请求参数错误", c)
		return
	}

	a, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Error(result.CodeNotFound, "文章不存在", c)
			return
		}
		result.Error(result.CodeServerError, "服务器内部错误", c)
		return
	}
	result.Success(toArticleResponse(a), c)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		result.Error(result.CodeParamError, "ID格式错误", c)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Error(result.CodeNotFound, "文章不存在", c)
			return
		}
		result.Error(result.CodeServerError, "服务器内部错误", c)
		return
	}
	result.Success(nil, c)
}
