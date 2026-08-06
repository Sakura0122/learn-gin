package user

import (
	"context"
	"errors"
	"learn-gin/internal/common/page"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrUsernameExists = errors.New("用户名已存在")
	userSortFields    = map[string]string{
		"id":         "id",
		"username":   "username",
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

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	count, err := gorm.G[User](s.db).Where("username = ?", req.Username).Count(ctx, "*")
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUsernameExists
	}

	u := User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
	}
	if err := gorm.G[User](s.db).Create(ctx, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	u, err := gorm.G[User](s.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) List(ctx context.Context, pageRequest page.Request) ([]User, int64, error) {
	orderBy, err := pageRequest.OrderBy(userSortFields, "id DESC")
	if err != nil {
		return nil, 0, err
	}

	total, err := gorm.G[User](s.db).Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}
	list, err := gorm.G[User](s.db).
		Offset(pageRequest.Offset()).
		Limit(pageRequest.PageSize).
		Order(orderBy).
		Find(ctx)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*User, error) {
	if _, err := s.GetByID(ctx, id); err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) > 0 {
		if _, err := gorm.G[User](s.db).
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
	_, err := gorm.G[User](s.db).Where("id = ?", id).Delete(ctx)
	return err
}
