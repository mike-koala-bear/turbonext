package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Message struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	RoomID    uint           `json:"room_id"`
	Room      Room           `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Room represents a chat room
type Room struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Messages  []Message      `json:"messages"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"` // Password should be hashed
}

var db *gorm.DB

func initDatabase() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read the environment variables for the DSN
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	// Construct the DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	// Connect to the database
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Automatically migrate the schema
	err = db.AutoMigrate(&Message{}, &Room{}, &User{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	log.Println("Database connection and migration successful")
}
