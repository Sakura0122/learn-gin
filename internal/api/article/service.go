package article

import (
	"context"
	"errors"

	"learn-gin/internal/api/user"
	"learn-gin/internal/common/page"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrUserNotFound   = errors.New("用户不存在")
	articleSortFields = map[string]string{
		"id":         "id",
		"user_id":    "user_id",
		"title":      "title",
		"status":     "status",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ctx context.Context, req CreateArticleRequest) (*Article, error) {
	count, err := gorm.G[user.User](s.db).Where("id = ?", req.UserID).Count(ctx, "*")
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrUserNotFound
	}

	a := Article{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}
	if err := gorm.G[Article](s.db).Create(ctx, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Article, error) {
	a, err := gorm.G[Article](s.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) List(ctx context.Context, userID uuid.UUID, pageRequest page.Request) ([]Article, int64, error) {
	orderBy, err := pageRequest.OrderBy(articleSortFields, "id DESC")
	if err != nil {
		return nil, 0, err
	}

	query := gorm.G[Article](s.db).Where("1 = 1")
	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}
	total, err := query.Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}
	list, err := query.
		Offset(pageRequest.Offset()).
		Limit(pageRequest.PageSize).
		Order(orderBy).
		Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateArticleRequest) (*Article, error) {
	if _, err := s.GetByID(ctx, id); err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) > 0 {
		if _, err := gorm.G[Article](s.db).
			Where("id = ?", id).
			Set(clause.Assignments(updates)).
			Update(ctx); err != nil {
			return nil, err
		}
	}
	return s.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	_, err := gorm.G[Article](s.db).Where("id = ?", id).Delete(ctx)
	return err
}
