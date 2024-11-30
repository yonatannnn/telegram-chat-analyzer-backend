// cmd/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"telegram-chat-analyzer/internal/delivery"
	"telegram-chat-analyzer/internal/repository"
	"telegram-chat-analyzer/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	mongoURI := os.Getenv("MONGO_URI")
	fmt.Println(mongoURI)
	dbName := os.Getenv("MONGO_DB")
	collection := os.Getenv("MONGO_COLLECTION")

	repo, err := repository.NewMongoRepository(mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize use case
	uc := usecase.NewMessageUsecase()

	// Set up Gin
	r := gin.Default()

	// Initialize handlers
	delivery.NewMessageHandler(r, uc, repo, collection)

	// Run server
	log.Println("Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
