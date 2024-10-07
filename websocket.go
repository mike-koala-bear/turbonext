package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Define a thread-safe map for clients per room
var clients = make(map[uint]map[*websocket.Conn]bool) // roomID -> clients
var clientsMutex = sync.RWMutex{}                     // Mutex to protect the clients map

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handle incoming WebSocket connections
func handleWebSocket(c *gin.Context) {
	roomIDParam := c.Param("roomID")
	roomIDUint, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}
	roomID := uint(roomIDUint)

	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	// Retrieve the username from the authenticated context
	username, exists := c.Get("username")
	if !exists {
		log.Printf("Username not found in context")
		ws.WriteJSON(gin.H{"error": "Unauthorized"})
		ws.Close()
		return
	}

	// Add the client to the room's client list
	clientsMutex.Lock()
	if clients[roomID] == nil {
		clients[roomID] = make(map[*websocket.Conn]bool)
	}
	clients[roomID][ws] = true
	clientsMutex.Unlock()
	log.Printf("Client %s connected to room %d", username.(string), roomID)

	defer func() {
		// Remove the client from the room's client list on disconnect
		clientsMutex.Lock()
		delete(clients[roomID], ws)
		clientsMutex.Unlock()
		log.Printf("Client %s disconnected from room %d", username.(string), roomID)
		ws.Close()
	}()

	// Listen for incoming messages
	for {
		var incoming struct {
			Content string `json:"content"`
		}

		// Read JSON message from WebSocket
		err := ws.ReadJSON(&incoming)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket closure: %v", err)
			}
			break
		}

		// Validate incoming message
		if incoming.Content == "" {
			ws.WriteJSON(gin.H{"error": "Message content cannot be empty"})
			continue
		}

		// Create a new Message instance
		newMessage := Message{
			Username:  username.(string),
			Content:   incoming.Content,
			CreatedAt: time.Now(),
			RoomID:    roomID,
		}

		// Save the message to the database
		if err := db.Create(&newMessage).Error; err != nil {
			log.Printf("Failed to save message: %v", err)
			ws.WriteJSON(gin.H{"error": "Failed to save message"})
			continue
		}

		// Broadcast the message to all clients in the same room
		broadcast <- newMessage
	}
}
