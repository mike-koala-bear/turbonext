package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

// Update the clients map to use uint as the key
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
	roomIDUint64, err := strconv.ParseUint(roomIDParam, 10, strconv.IntSize)
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Ensure that the parsed roomID fits into a uint
	maxUint := uint64(^uint(0))
	if roomIDUint64 > maxUint {
		log.Printf("Room ID out of bounds: %v", roomIDUint64)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID out of bounds"})
		return
	}
	roomID := uint(roomIDUint64)

	// Extract token from cookies
	tokenString, err := c.Cookie("jwt_token")
	if err != nil || tokenString == "" {
		log.Printf("Missing or invalid token for WebSocket connection")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		return
	}

	// Validate JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		log.Printf("Invalid or expired token for WebSocket connection: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	username := claims.Username

	// Add the client to the room's client list
	clientsMutex.Lock()
	if clients[roomID] == nil {
		clients[roomID] = make(map[*websocket.Conn]bool)
	}
	clients[roomID][ws] = true
	clientsMutex.Unlock()
	log.Printf("Client %s connected to room %d via WebSocket", username, roomID)

	defer func() {
		// Remove the client from the room's client list on disconnect
		clientsMutex.Lock()
		delete(clients[roomID], ws)
		clientsMutex.Unlock()
		log.Printf("Client %s disconnected from room %d", username, roomID)
		ws.Close()
	}()

	// Since messages are sent via HTTP POST, WebSocket is only for receiving
	for {
		// Listen for close messages or pings/pongs to keep the connection alive
		if _, _, err := ws.NextReader(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket closure for client %s in room %d: %v", username, roomID, err)
			}
			break
		}
	}
}
