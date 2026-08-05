package user

import (
	"errors"

	"github.com/google/uuid"
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

func (s *Service) Create(req CreateUserRequest) (*User, error) {
	var count int64
	if err := s.db.Model(&User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
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
	if err := s.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) GetByID(id uuid.UUID) (*User, error) {
	var u User
	if err := s.db.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) List(page, pageSize int) ([]User, int64, error) {
	var (
		list  []User
		total int64
	)
	if err := s.db.Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := s.db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *Service) Update(id uuid.UUID, req UpdateUserRequest) (*User, error) {
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
		if err := s.db.Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetByID(id)
}

func (s *Service) Delete(id uuid.UUID) error {
	if _, err := s.GetByID(id); err != nil {
		return err
	}
	return s.db.Delete(&User{}, "id = ?", id).Error
}
