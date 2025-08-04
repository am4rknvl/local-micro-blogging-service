package handlers

import (
	"log"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/gofiber/fiber/v2"
)

func SendFriendRequest(c *fiber.Ctx) error {
	senderID := c.Get("X-User-ID") // adjust if you use middleware Locals("userID")

	var payload struct {
		ReceiverID string `json:"receiver_id"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	err := db.DB.Exec(`INSERT INTO friend_requests (sender_id, receiver_id) VALUES (?, ?)`, senderID, payload.ReceiverID).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not send request")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Friend request sent"})
}

func RespondToFriendRequest(c *fiber.Ctx) error {
	var payload struct {
		RequestID int    `json:"request_id"`
		Action    string `json:"action"` // accept or reject
	}

	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
	}

	var friendRequest struct {
		SenderID   string `gorm:"column:sender_id"`
		ReceiverID string `gorm:"column:receiver_id"`
	}

	err := db.DB.Table("friend_requests").Where("id = ?", payload.RequestID).First(&friendRequest).Error
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Request not found")
	}

	senderID := friendRequest.SenderID
	receiverID := friendRequest.ReceiverID

	if payload.Action == "accept" {
		_ = db.DB.Exec(`UPDATE friend_requests SET status='accepted' WHERE id=?`, payload.RequestID)
		_ = db.DB.Exec(`INSERT INTO friends (user1_id, user2_id) VALUES (?, ?)`, senderID, receiverID)

		// Auto-follow both directions
		_ = db.DB.Exec(`
			INSERT INTO follows (id, follower_id, followee_id, created_at)
			SELECT gen_random_uuid(), ?, ?, NOW()
			WHERE NOT EXISTS (
				SELECT 1 FROM follows WHERE follower_id=? AND followee_id=?
			)
		`, senderID, receiverID, senderID, receiverID)

		_ = db.DB.Exec(`
			INSERT INTO follows (id, follower_id, followee_id, created_at)
			SELECT gen_random_uuid(), ?, ?, NOW()
			WHERE NOT EXISTS (
				SELECT 1 FROM follows WHERE follower_id=? AND followee_id=?
			)
		`, receiverID, senderID, receiverID, senderID)

		return c.JSON(fiber.Map{"message": "Friend request accepted"})
	}

	// Rejected path, just update status
	_ = db.DB.Exec(`UPDATE friend_requests SET status='rejected' WHERE id=?`, payload.RequestID)
	return c.JSON(fiber.Map{"message": "Friend request rejected"})
}

func GetFriendRequests(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	type FriendRequestInfo struct {
		ID          int    `json:"id"`
		SenderID    string `json:"sender_id"`
		Username    string `json:"username"`
		Avatar      string `json:"avatar"`
		RequestedAt string `json:"created_at"`
	}

	var requests []FriendRequestInfo

	err := db.DB.Raw(`
		SELECT fr.id, u.id as sender_id, u.username, u.avatar, fr.created_at
		FROM friend_requests fr
		JOIN users u ON u.id = fr.sender_id
		WHERE fr.receiver_id = ? AND fr.status = 'pending'
	`, userID).Scan(&requests).Error

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load friend requests")
	}

	return c.JSON(requests)
}

func GetFriendTree(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	query := `
		SELECT u.id, u.username, u.avatar
		FROM follows f1
		JOIN follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
		JOIN users u ON u.id = f1.following_id
		WHERE f1.follower_id = ?
	`

	rows, err := db.DB.Raw(query, userID).Rows()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	defer rows.Close()

	friendTree := []fiber.Map{}

	for rows.Next() {
		var friendID, friendUsername, friendAvatar string
		if err := rows.Scan(&friendID, &friendUsername, &friendAvatar); err != nil {
			continue
		}

		mutualRows, err := db.DB.Raw(query, friendID).Rows()
		if err != nil {
			log.Println("mutuals query error:", err)
			continue
		}

		var mutuals []fiber.Map
		for mutualRows.Next() {
			var mid, mname, mavatar string
			if err := mutualRows.Scan(&mid, &mname, &mavatar); err == nil && mid != userID {
				mutuals = append(mutuals, fiber.Map{
					"id":       mid,
					"username": mname,
					"avatar":   mavatar,
				})
			}
		}
		mutualRows.Close()

		friendTree = append(friendTree, fiber.Map{
			"id":       friendID,
			"username": friendUsername,
			"avatar":   friendAvatar,
			"mutuals":  mutuals,
		})
	}

	return c.JSON(fiber.Map{
		"user":    userID,
		"friends": friendTree,
	})
}

func GetFriends(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	query := `
		SELECT u.id, u.username, u.avatar
		FROM follows f1
		JOIN follows f2 ON f1.follower_id = f2.following_id AND f1.following_id = f2.follower_id
		JOIN users u ON u.id = f1.following_id
		WHERE f1.follower_id = ?
	`

	rows, err := db.DB.Raw(query, userID).Rows()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	defer rows.Close()

	var friends []fiber.Map
	for rows.Next() {
		var id, username, avatar string
		if err := rows.Scan(&id, &username, &avatar); err == nil {
			friends = append(friends, fiber.Map{
				"id":       id,
				"username": username,
				"avatar":   avatar,
			})
		}
	}
	return c.JSON(fiber.Map{"friends": friends})
}
