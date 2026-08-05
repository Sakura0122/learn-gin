package article

import (
	"github.com/google/uuid"
	"learn-gin/internal/infra/database"
)

type Article struct {
	database.BaseTable
	UserID  uuid.UUID `gorm:"type:char(36);not null;index" json:"user_id"`
	Title   string    `gorm:"size:128;not null" json:"title"`
	Content string    `gorm:"type:text" json:"content"`
	Status  int       `gorm:"not null;default:1" json:"status"`
}

func (Article) TableName() string {
	return "article"
}
