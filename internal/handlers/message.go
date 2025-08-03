package handlers

import (
	"time"
	"database/sql"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func SendMessage(c *fiber.Ctx) error {
	senderID := c.Locals("userID").(string)
	convoID := c.Params("id")

	var msg models.Message
	if err := c.BodyParser(&msg); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	msg.SenderID = senderID
	msg.ConversationID = convoID
	msg.CreatedAt = time.Now()

	
	if err := db.DB.Create(&msg).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send message"})
	}

	return c.Status(201).JSON(msg)
}

func GetMessages(c *fiber.Ctx) error {
	convoID := c.Params("id")

	var messages []models.Message
	if err := db.DB.Where("conversation_id = ?", convoID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}

	return c.JSON(messages)
}


func SaveMessage(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		messageID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid message ID"})
		}

		_, err = db.Exec("UPDATE messages SET is_saved = TRUE WHERE id = $1", messageID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save message"})
		}

		return c.SendStatus(fiber.StatusNoContent) // 204 No Content
	}
}
