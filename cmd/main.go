func main() {
    // Try to load .env file, but don't fail if it doesn't exist
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found; falling back to environment variables.")
    }

    // Read environment variables
    mongoURI := os.Getenv("MONGO_URI")
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

    // Initialize handlers
    delivery.NewMessageHandler(r, uc, repo, collection)

    // Run server
    log.Println("Server running on port 8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
