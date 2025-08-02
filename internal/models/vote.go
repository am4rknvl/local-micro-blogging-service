package models

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Value     int8      `gorm:"not null"` // +1 for upvote, -1 for downvote, 0 for no vote
	CreatedAt time.Time
	UpdatedAt time.Time
}
