// controllers/ws_controller.go

package controllers

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// ConnectionManager handles active WebSocket connections
type ConnectionManager struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// NewConnectionManager initializes a new WebSocket connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

var connManager = NewConnectionManager()

// WebSocketHandler handles new WebSocket connections
func WebSocketHandler(c *websocket.Conn) {
	connManager.mu.Lock()
	connManager.clients[c] = true
	connManager.mu.Unlock()

	defer func() {
		connManager.mu.Lock()
		delete(connManager.clients, c)
		connManager.mu.Unlock()
		c.Close()
	}()

	for {
		// WebSocket requires continuous message reads to keep connection open
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}
}

// BroadcastUpdate sends a message to all connected WebSocket clients
func BroadcastUpdate(message string) {
	connManager.mu.Lock()
	defer connManager.mu.Unlock()

	for client := range connManager.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
			client.Close()
			delete(connManager.clients, client)
		}
	}
}
