package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `json:"-"`
	Avatar    string    `json:"avatar,omitempty"`
	Bio       string    `gorm:"size:500" json:"bio,omitempty"`
	Website   string    `gorm:"size:255" json:"website,omitempty"`
	Location  string    `gorm:"size:100" json:"location,omitempty"`
	IsPrivate bool      `gorm:"default:false" json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Other fields...
}
