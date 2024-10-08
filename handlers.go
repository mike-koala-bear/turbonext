package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func getTokenHandler(c *gin.Context) {
	// Retrieve JWT token from cookie
	tokenString, err := c.Cookie("jwt_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func checkAuthHandler(c *gin.Context) {
	// Retrieve JWT token from HttpOnly cookie
	tokenString, err := c.Cookie("jwt_token")
	if err != nil || tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Validate JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Return success if the token is valid
	c.JSON(http.StatusOK, gin.H{"message": "Authenticated"})
}

// Create a new room
func createRoomHandler(c *gin.Context) {
	var newRoom struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	room := Room{
		Name:      newRoom.Name,
		CreatedAt: time.Now(),
	}

	if err := db.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// Get all rooms
func getRoomsHandler(c *gin.Context) {
	var rooms []Room
	if err := db.Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rooms"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

// Get messages in a room
func getMessagesHandler(c *gin.Context) {
	roomIDParam := c.Param("roomID")
	roomIDUint64, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Calculate max uint value based on platform
	maxUint := ^uint(0)
	if roomIDUint64 > uint64(maxUint) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID out of bounds"})
		return
	}
	roomID := uint(roomIDUint64)

	var messages []Message
	if err := db.Where("room_id = ?", roomID).Order("created_at desc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// postMessageHandler handles incoming messages via HTTP POST
func postMessageHandler(c *gin.Context) {
	roomIDParam := c.Param("roomID")
	roomIDUint64, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Calculate max uint value based on platform
	maxUint := ^uint(0)
	if roomIDUint64 > uint64(maxUint) {
		log.Printf("Room ID out of bounds: %v", roomIDUint64)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID out of bounds"})
		return
	}
	roomID := uint(roomIDUint64)

	// Extract token from cookies
	tokenString, err := c.Cookie("jwt_token")
	if err != nil || tokenString == "" {
		log.Printf("Missing or invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		return
	}

	// Validate JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		log.Printf("Invalid or expired token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Bind JSON input
	var incoming struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&incoming); err != nil {
		log.Printf("Message binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate content
	if incoming.Content == "" {
		log.Printf("Received empty message from %s in room %d", claims.Username, roomID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message content cannot be empty"})
		return
	}

	// Create a new Message instance
	newMessage := Message{
		Username:  claims.Username,
		Content:   incoming.Content,
		CreatedAt: time.Now(),
		RoomID:    roomID,
	}

	// Save the message to the database
	if err := db.Create(&newMessage).Error; err != nil {
		log.Printf("Failed to save message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	log.Printf("Message from %s in room %d saved to database with ID %d", claims.Username, roomID, newMessage.ID)

	// Broadcast the message to all clients in the same room
	Broadcast <- newMessage

	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully", "data": newMessage})
}
