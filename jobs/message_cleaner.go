package jobs

import (
	"log"
	"time"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
)

func StartMessageCleanupJob() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // runs every hour
		for range ticker.C {
			deleteExpiredMessages()
		}
	}()
}

func deleteExpiredMessages() {
	query := `
		DELETE FROM messages 
		WHERE is_saved = false AND created_at < NOW() - INTERVAL '24 HOURS';
	`
	err := db.DB.Exec(query).Error
	if err != nil {
		log.Println("Failed to delete expired messages:", err)
	} else {
		log.Println("Ephemeral messages cleaned up.")
	}
}
