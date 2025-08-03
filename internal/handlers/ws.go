package handlers

import (
	"encoding/json"
	"log"

	db "github.com/am4rknvl/local-micro-blogging-service.git/internal/database"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/models"
	"github.com/am4rknvl/local-micro-blogging-service.git/internal/ws"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type IncomingMessage struct {
	ConversationID string `json:"conversation_id"`
	Content        string `json:"content"`
}

func WebSocketHandler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		userID := c.Locals("userID").(string) // from JWT middleware
		convoID := c.Params("conversationID") // passed in URL

		client := &ws.Client{
			Conn:           c,
			UserID:         userID,
			ConversationID: convoID,
		}

		ws.ManagerInstance.Register <- client
		defer func() {
			ws.ManagerInstance.Unregister <- client
			c.Close()
		}()

		for {
			_, data, err := c.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				break
			}

			var incoming IncomingMessage
			if err := json.Unmarshal(data, &incoming); err != nil {
				continue
			}

			// Save to DB
			msg := &models.Message{
				SenderID:       userID,
				ConversationID: convoID,
				Content:        incoming.Content,
				IsSaved:        false,
			}
			// Save to DB using GORM
			if err := db.DB.Create(msg).Error; err != nil {
				log.Println("DB save failed:", err)
				continue
			}
			saved := msg

			// Broadcast back to all in this convo
			outgoing, _ := json.Marshal(saved)
			ws.ManagerInstance.Broadcast <- ws.MessagePayload{
				ConversationID: convoID,
				Data:           outgoing,
			}
		}
	})
}
