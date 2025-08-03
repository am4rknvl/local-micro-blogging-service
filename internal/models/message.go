package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConversationID string    `json:"conversation_id" gorm:"not null"`
	SenderID       string    `json:"sender_id" gorm:"not null"`
	Content        string    `json:"content" gorm:"not null"`
	Type           string    `json:"type" gorm:"default:'text'"`
	IsSaved        bool      `json:"is_saved" gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
