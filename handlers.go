package main

import (
	"log"
	"math"
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

// Get messages in a room along with the room name
func getMessagesHandler(c *gin.Context) {
	roomIDParam := c.Param("roomID")
	roomIDUint64, err := strconv.ParseUint(roomIDParam, 10, strconv.IntSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Perform upper bound check using constants from the math package
	var maxUint64 uint64
	if strconv.IntSize == 32 {
		maxUint64 = uint64(math.MaxUint32)
	} else if strconv.IntSize == 64 {
		maxUint64 = uint64(math.MaxUint64)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unsupported platform"})
		return
	}

	if roomIDUint64 > maxUint64 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID out of bounds"})
		return
	}

	roomID := uint(roomIDUint64)

	// Fetch the room name
	var room Room
	if err := db.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	// Fetch messages in the room
	var messages []Message
	if err := db.Where("room_id = ?", roomID).Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	// Return both room name and messages
	c.JSON(http.StatusOK, gin.H{
		"room_name": room.Name,
		"messages":  messages,
	})
}

// postMessageHandler handles incoming messages via HTTP POST
func postMessageHandler(c *gin.Context) {
	roomIDParam := c.Param("roomID")
	roomIDUint64, err := strconv.ParseUint(roomIDParam, 10, strconv.IntSize)
	if err != nil {
		log.Printf("Invalid room ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	// Perform upper bound check using constants from the math package
	var maxUint64 uint64
	if strconv.IntSize == 32 {
		maxUint64 = uint64(math.MaxUint32)
	} else if strconv.IntSize == 64 {
		maxUint64 = uint64(math.MaxUint64)
	} else {
		log.Printf("Unsupported platform")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unsupported platform"})
		return
	}

	if roomIDUint64 > maxUint64 {
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
