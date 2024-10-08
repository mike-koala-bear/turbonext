package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	initDatabase()

	// Create a new Gin router
	r := gin.Default()

	// Set trusted proxies for development environment (localhost only)
	proxyerr := r.SetTrustedProxies([]string{"127.0.0.1"})
	if proxyerr != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	// Enable CORS
	r.Use(cors.Default())

	// Public routes
	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)
	r.POST("/logout", logoutHandler)

	// Protected routes
	protected := r.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/rooms", getRoomsHandler)
		protected.POST("/rooms", createRoomHandler)
		protected.GET("/rooms/:roomID/messages", getMessagesHandler)
		protected.POST("/rooms/:roomID/messages", postMessageHandler)
		protected.GET("/ws/:roomID", handleWebSocket)
		protected.GET("/get-token", getTokenHandler)
		protected.GET("/check-auth", checkAuthHandler)
	}

	// Start message broadcasting in the background
	go handleBroadcast()

	// Start the server
	r.Run(":8080")
}
