package article

import (
	"learn-gin/internal/infra/database"
)

type ArticleResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toArticleResponse(a *database.Article) ArticleResponse {
	return ArticleResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Title:     a.Title,
		Content:   a.Content,
		Status:    a.Status,
		CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: a.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func toArticleResponseList(list []database.Article) []ArticleResponse {
	res := make([]ArticleResponse, 0, len(list))
	for i := range list {
		res = append(res, toArticleResponse(&list[i]))
	}
	return res
}
