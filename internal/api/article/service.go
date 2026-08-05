package article

import (
	"errors"

	"learn-gin/internal/infra/database"

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

func (s *Service) Create(req CreateArticleRequest) (*database.Article, error) {
	var count int64
	if err := s.db.Model(&database.User{}).Where("id = ?", req.UserID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, ErrUserNotFound
	}

	a := database.Article{
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
	}
	if err := s.db.Create(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) GetByID(id uint) (*database.Article, error) {
	var a database.Article
	if err := s.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *Service) List(userID uint, page, pageSize int) ([]database.Article, int64, error) {
	var (
		list  []database.Article
		total int64
	)
	query := s.db.Model(&database.Article{})
	if userID > 0 {
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

func (s *Service) Update(id uint, req UpdateArticleRequest) (*database.Article, error) {
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
		if err := s.db.Model(&database.Article{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.db.Delete(&database.Article{}, id).Error
}
