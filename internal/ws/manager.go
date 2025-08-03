package ws

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn           *websocket.Conn
	UserID         string
	ConversationID string
}

type Manager struct {
	mu         sync.RWMutex
	clients    map[string][]*Client // conversationID -> []*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan MessagePayload
}

type MessagePayload struct {
	ConversationID string
	Data           []byte
}

// Global manager instance
var ManagerInstance = NewManager()

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[string][]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan MessagePayload),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.Register:
			m.mu.Lock()
			m.clients[client.ConversationID] = append(m.clients[client.ConversationID], client)
			m.mu.Unlock()

		case client := <-m.Unregister:
			m.mu.Lock()
			convoClients := m.clients[client.ConversationID]
			for i, c := range convoClients {
				if c == client {
					m.clients[client.ConversationID] = append(convoClients[:i], convoClients[i+1:]...)
					break
				}
			}
			m.mu.Unlock()

		case payload := <-m.Broadcast:
			m.mu.RLock()
			convoClients := m.clients[payload.ConversationID]
			m.mu.RUnlock()
			for _, c := range convoClients {
				if err := c.Conn.WriteMessage(websocket.TextMessage, payload.Data); err != nil {
					log.Println("WS write error:", err)
				}
			}
		}
	}
}
