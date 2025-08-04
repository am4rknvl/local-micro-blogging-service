package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	User struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Avatar   string    `json:"avatar"`
	} `json:"user"`
}


