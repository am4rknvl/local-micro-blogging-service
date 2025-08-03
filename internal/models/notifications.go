package models

import (
	"time"
)

type Notification struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`       // Who receives it
	ActorID     int64     `json:"actor_id"`      // Who triggered it
	Type        string    `json:"type"`          // e.g., "comment", "follow", "message", etc.
	EntityID    int64     `json:"entity_id"`     // ID of related entity (post, message, etc.)
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}
