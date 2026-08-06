package user

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
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(result.CodeParamError, "请求参数错误", c)
		return
	}

	u, err := h.service.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUsernameExists):
			result.Error(result.CodeError, "用户名已存在", c)
		default:
			result.Error(result.CodeServerError, "服务器内部错误", c)
		}
		return
	}
	result.Success(toUserResponse(u), c)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		result.Error(result.CodeParamError, "ID格式错误", c)
		return
	}

	u, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Error(result.CodeNotFound, "用户不存在", c)
			return
		}
		result.Error(result.CodeServerError, "服务器内部错误", c)
		return
	}
	result.Success(toUserResponse(u), c)
}

func (h *Handler) List(c *gin.Context) {
	pageRequest, err := page.Parse(c)
	if err != nil {
		result.Error(result.CodeParamError, "分页参数错误", c)
		return
	}

	list, total, err := h.service.List(pageRequest)
	if err != nil {
		if errors.Is(err, page.ErrUnsupportedSortField) {
			result.Error(result.CodeParamError, err.Error(), c)
			return
		}
		result.Error(result.CodeServerError, "服务器内部错误", c)
		return
	}
	result.Success(page.NewResult(pageRequest, total, toUserResponseList(list)), c)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		result.Error(result.CodeParamError, "ID格式错误", c)
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(result.CodeParamError, "请求参数错误", c)
		return
	}

	u, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Error(result.CodeNotFound, "用户不存在", c)
			return
		}
		result.Error(result.CodeServerError, "服务器内部错误", c)
		return
	}
	result.Success(toUserResponse(u), c)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		result.Error(result.CodeParamError, "ID格式错误", c)
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Error(result.CodeNotFound, "用户不存在", c)
			return
		}
		result.Error(result.CodeServerError, "服务器内部错误", c)
		return
	}
	result.Success(nil, c)
}
