package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	dbName     = "OrgDB"
)

// InitDB initializes the MongoDB database connection.
func InitDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	// Set the global client variable for use throughout the application
	SetMongoClient(client)
}

// SetMongoClient sets the MongoDB client.
func SetMongoClient(c *mongo.Client) {
	client = c
}

// GetMongoClient retrieves the MongoDB client.
func GetMongoClient() *mongo.Client {
	return client
}
func GetUsersCollection() *mongo.Collection {
	return client.Database(dbName).Collection("users")
}
