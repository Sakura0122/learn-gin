package article

type CreateArticleRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type UpdateArticleRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  *int    `json:"status"`
}
