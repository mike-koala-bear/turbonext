package main

import (
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
	roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var messages []Message
	if err := db.Where("room_id = ?", roomID).Order("created_at desc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// Post a message to a room
func postMessageHandler(c *gin.Context) {
	roomIDParam := c.Param("roomID")
	roomID, err := strconv.ParseUint(roomIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	var newMessage struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	message := Message{
		Username:  username.(string),
		Content:   newMessage.Content,
		CreatedAt: time.Now(),
		RoomID:    uint(roomID),
	}

	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post message"})
		return
	}

	// Broadcast the message to WebSocket clients in the room
	broadcast <- message
	c.JSON(http.StatusCreated, message)
}
