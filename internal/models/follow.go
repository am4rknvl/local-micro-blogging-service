package models

import (
	"time"
	"github.com/google/uuid"
)


type Follow struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	FollowerID uuid.UUID `gorm:"type:uuid;not null;index"` // user who follows
	FolloweeID uuid.UUID `gorm:"type:uuid;not null;index"` // user being followed
	CreatedAt time.Time
}
