package user

import (
	"learn-gin/internal/infra/database"
)

type User struct {
	database.BaseTable
	Username string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:128;not null" json:"-"`
	Nickname string `gorm:"size:64;not null;default:''" json:"nickname"`
	Email    string `gorm:"size:128;not null;default:''" json:"email"`
	Phone    string `gorm:"size:20;not null;default:''" json:"phone"`
	Status   int    `gorm:"not null;default:1" json:"status"`
}

func (User) TableName() string {
	return "user"
}
