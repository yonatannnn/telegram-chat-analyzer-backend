package main

import (
	"log"
	"os"
	"telegram-chat-analyzer/internal/delivery"
	"telegram-chat-analyzer/internal/repository"
	"telegram-chat-analyzer/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; falling back to environment variables.")
	}

	// Read environment variables
	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		log.Fatal("MONGO_DB is not set")
	}
	collection := os.Getenv("MONGO_COLLECTION")
	if collection == "" {
		log.Fatal("MONGO_COLLECTION is not set")
	}

	// Initialize MongoDB repository
	repo, err := repository.NewMongoRepository(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize use case
	uc := usecase.NewMessageUsecase()

	// Set up Gin
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize handlers
	delivery.NewMessageHandler(r, uc, repo, collection)

	// Run server
	log.Println("Server running on port 8080")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
