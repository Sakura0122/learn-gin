package article

import "github.com/google/uuid"

type CreateArticleRequest struct {
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Title   string    `json:"title" binding:"required"`
	Content string    `json:"content"`
}

type UpdateArticleRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  *int    `json:"status"`
}
