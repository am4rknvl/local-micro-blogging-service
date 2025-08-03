package jobs

import (
	"database/sql"
	"log"
	"time"
)

func StartAutoDeleteJob(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			if err := deleteOldMessages(db); err != nil {
				log.Println("Error deleting old messages:", err)
			} else {
				log.Println("🧹 Old unsaved messages deleted.")
			}
		}
	}()
}

func deleteOldMessages(db *sql.DB) error {
	_, err := db.Exec(`
		DELETE FROM messages
		WHERE is_saved = FALSE
		AND created_at < NOW() - INTERVAL '24 HOURS'`)
	return err
}

func DeleteOldMessages(db *sql.DB) error {
	return deleteOldMessages(db)
}
