package user

import (
	"errors"

	"learn-gin/internal/infra/database"

	"gorm.io/gorm"
)

var (
	ErrUsernameExists = errors.New("username already exists")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(req CreateUserRequest) (*database.User, error) {
	var count int64
	if err := s.db.Model(&database.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUsernameExists
	}

	u := database.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
	}
	if err := s.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) GetByID(id uint) (*database.User, error) {
	var u database.User
	if err := s.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) List(page, pageSize int) ([]database.User, int64, error) {
	var (
		list  []database.User
		total int64
	)
	if err := s.db.Model(&database.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := s.db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *Service) Update(id uint, req UpdateUserRequest) (*database.User, error) {
	if _, err := s.GetByID(id); err != nil {
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
		if err := s.db.Model(&database.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetByID(id)
}

func (s *Service) Delete(id uint) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.db.Delete(&database.User{}, id).Error
}
