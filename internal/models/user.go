package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username  string    `gorm:"uniqueIndex"`
	Email     string    `gorm:"uniqueIndex"`
	Password  string
	Avatar    string
	Blocked   bool
	CreatedAt time.Time
}
