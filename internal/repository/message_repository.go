// internal/repository/mongo_repository.go
package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository interface {
	SaveProcessedData(ctx context.Context, collection string, data interface{}) error
}

type mongoRepository struct {
	client *mongo.Client
	dbName string
}

func NewMongoRepository(connectionString, dbName string) (MongoRepository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &mongoRepository{client: client, dbName: dbName}, nil
}

func (r *mongoRepository) SaveProcessedData(ctx context.Context, collection string, data interface{}) error {
	coll := r.client.Database(r.dbName).Collection(collection)

	_, err := coll.InsertOne(ctx, data)
	if err != nil {
		log.Printf("Failed to save data to MongoDB: %v", err)
		return err
	}

	log.Println("Data successfully saved to MongoDB!")
	return nil
}
