package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func getMessagesHandler(c *gin.Context) {
	var messages []Message
	if err := db.Order("created_at desc").Find(&messages).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve messages"})
		return
	}
	c.JSON(200, messages)
}

func postMessageHandler(c *gin.Context) {
	var newMessage Message
	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	newMessage.CreatedAt = time.Now()
	if err := db.Create(&newMessage).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to post message"})
		return
	}

	// Broadcast the message to WebSocket clients
	broadcast <- newMessage
	c.JSON(201, newMessage)
}

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
