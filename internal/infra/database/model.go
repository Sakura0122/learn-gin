package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:128;not null" json:"-"`
	Nickname  string         `gorm:"size:64;not null;default:''" json:"nickname"`
	Email     string         `gorm:"size:128;not null;default:''" json:"email"`
	Phone     string         `gorm:"size:20;not null;default:''" json:"phone"`
	Status    int            `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "user"
}

type Article struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Title     string         `gorm:"size:128;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Status    int            `gorm:"not null;default:1" json:"status"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Article) TableName() string {
	return "article"
}
