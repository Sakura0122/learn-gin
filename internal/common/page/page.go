package page

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrUnsupportedSortField = errors.New("排序字段不支持")

type Request struct {
	CurrentPage int    `form:"page,default=1" binding:"gte=1"`
	PageSize    int    `form:"page_size,default=10" binding:"gte=1"`
	SortField   string `form:"sort_field"`
	IsAsc       bool   `form:"is_asc,default=true"`
}

func Parse(c *gin.Context) (Request, error) {
	var request Request
	if err := c.ShouldBindQuery(&request); err != nil {
		return Request{}, err
	}
	return request, nil
}

func (r Request) Offset() int {
	return (r.CurrentPage - 1) * r.PageSize
}

func (r Request) OrderBy(allowedSortFields map[string]string, defaultOrder string) (string, error) {
	if r.SortField == "" {
		return defaultOrder, nil
	}

	column, ok := allowedSortFields[r.SortField]
	if !ok {
		return "", ErrUnsupportedSortField
	}

	direction := "DESC"
	if r.IsAsc {
		direction = "ASC"
	}
	return column + " " + direction, nil
}

type Result[T any] struct {
	Total       int64 `json:"total"`
	PageCount   int64 `json:"page_count"`
	Items       []T   `json:"list"`
	CurrentPage int   `json:"page"`
	PageSize    int   `json:"page_size"`
}

func NewResult[T any](request Request, total int64, items []T) Result[T] {
	pageSize := int64(request.PageSize)
	return Result[T]{
		Total:       total,
		PageCount:   (total + pageSize - 1) / pageSize,
		Items:       items,
		CurrentPage: request.CurrentPage,
		PageSize:    request.PageSize,
	}
}
