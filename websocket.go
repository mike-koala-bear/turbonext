package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

// upgrader upgrades the HTTP connection to a WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func handleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}

	// Add the client to the `clients` map
	clients[ws] = true
	log.Printf("Client connected: %v", ws.RemoteAddr())

	defer func() {
		// Remove the client from the `clients` map on disconnect
		delete(clients, ws)
		log.Printf("Client disconnected: %v", ws.RemoteAddr())
		ws.Close()
	}()

	var username string // Store the username

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var messageData map[string]string
		if err := json.Unmarshal(msg, &messageData); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Token Authentication (handle JWT here)
		if tokenStr, ok := messageData["token"]; ok {
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				ws.WriteMessage(websocket.TextMessage, []byte("Invalid token"))
				return
			}

			username = claims.Username
			log.Printf("User authenticated via WebSocket: %s", username)
			continue
		}

		// Process chat messages
		if content, ok := messageData["content"]; ok {
			if username == "" {
				ws.WriteMessage(websocket.TextMessage, []byte("User not authenticated"))
				continue
			}

			newMessage := Message{
				Username:  username,
				Content:   content,
				CreatedAt: time.Now(),
			}

			// Save the message to the database
			if err := db.Create(&newMessage).Error; err != nil {
				log.Printf("Failed to save message: %v", err)
			}

			// Broadcast the message
			broadcast <- newMessage
		}
	}
}
