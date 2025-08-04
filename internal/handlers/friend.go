package handlers

import (
	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/gofiber/fiber/v2"
)

func SendFriendRequest(c *fiber.Ctx) error {
	senderID := c.Get("X-User-ID")

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
		// Accept the request
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

	// If rejected, just update status — NO follow logic
	_ = db.DB.Exec(`UPDATE friend_requests SET status='rejected' WHERE id=?`, payload.RequestID)
	return c.JSON(fiber.Map{"message": "Friend request rejected"})
}


func GetFriendRequests(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	type FriendRequestInfo struct {
		ID         int    `json:"id"`
		SenderID   string `json:"sender_id"`
		Username   string `json:"username"`
		Avatar     string `json:"avatar"` // optional
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

	type UserNode struct {
		ID       string     `json:"id"`
		Username string     `json:"username"`
		Avatar   string     `json:"avatar"`
		Friends  []UserNode `json:"friends,omitempty"` // recursive 🌳
	}

	var direct []UserNode
	err := db.DB.Raw(`
		SELECT u.id, u.username, u.avatar
		FROM friends f
		JOIN users u ON (u.id = f.user1_id OR u.id = f.user2_id)
		WHERE (f.user1_id = ? OR f.user2_id = ?) AND u.id != ?
	`, userID, userID, userID).Scan(&direct).Error

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get friends")
	}

	// Build tree manually (second-degree layer only to keep sane)
	for i := range direct {
		friendID := direct[i].ID

		var second []UserNode
		_ = db.DB.Raw(`
			SELECT u.id, u.username, u.avatar
			FROM friends f
			JOIN users u ON (u.id = f.user1_id OR u.id = f.user2_id)
			WHERE (f.user1_id = ? OR f.user2_id = ?) AND u.id != ?
		`, friendID, friendID, friendID).Scan(&second)

		direct[i].Friends = second
	}

	return c.JSON(direct)
}
