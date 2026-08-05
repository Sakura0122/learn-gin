package article

import (
	"errors"

	"learn-gin/internal/api/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(req CreateArticleRequest) (*Article, error) {
	var count int64
	if err := s.db.Model(&user.User{}).Where("id = ?", req.UserID).Count(&count).Error; err != nil {
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
	if err := s.db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) GetByID(id uuid.UUID) (*Article, error) {
	var a Article
	if err := s.db.First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) List(userID uuid.UUID, page, pageSize int) ([]Article, int64, error) {
	var (
		list  []Article
		total int64
	)
	query := s.db.Model(&Article{})
	if userID != uuid.Nil {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *Service) Update(id uuid.UUID, req UpdateArticleRequest) (*Article, error) {
	if _, err := s.GetByID(id); err != nil {
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
		if err := s.db.Model(&Article{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetByID(id)
}

func (s *Service) Delete(id uuid.UUID) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.db.Delete(&Article{}, "id = ?", id).Error
}
