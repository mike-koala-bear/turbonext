package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var db *gorm.DB
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func initDatabase() {
	var err error
	dsn := "host=localhost user=turbonext password=turbonext_password dbname=turbonext port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	err = db.AutoMigrate(&Message{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	log.Println("Database connection and migration successful")
}

func main() {
	initDatabase()
	go handleMessages()

	r := gin.Default()

	// Set trusted proxies for development environment (localhost only)
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatal("Failed to set trusted proxies:", err)
	}

	// Enable CORS for all origins (adjust as needed for production)
	r.Use(cors.Default())

	// Routes
	r.GET("/messages", getMessagesHandler)
	r.POST("/messages", postMessageHandler)
	r.GET("/ws", handleWebSocket)

	r.Run(":8080")
}

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

	broadcast <- newMessage
	c.JSON(201, newMessage)
}

func handleWebSocket(c *gin.Context) {
    ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Failed to upgrade WebSocket: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open WebSocket connection"})
        return
    }

    log.Println("WebSocket upgrade successful")

    // Defer closure and cleanup when the connection is closed
    defer func() {
        log.Println("Closing WebSocket connection")
        delete(clients, ws)
        ws.Close()
    }()

    // Add WebSocket connection to the map of clients
    clients[ws] = true

    // Keep reading messages until an error occurs or connection closes
    for {
        _, _, err := ws.ReadMessage()
        if err != nil {
            // Log the error for debugging
            log.Printf("WebSocket read error: %v", err)
            delete(clients, ws)
            break
        }

        // If the connection is valid and no read error occurs, proceed with broadcasting
        var msg Message
        if err := db.Order("created_at desc").First(&msg).Error; err != nil {
            log.Printf("Failed to retrieve message: %v", err)
            continue
        }

        // Broadcast the message to all connected clients
        broadcast <- msg
    }
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
