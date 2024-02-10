package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	dbName string
)

// InitDB initializes the MongoDB database connection.
func InitDB(dataBaseName string,uri string) {
	dbName = dataBaseName
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
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

	SetMongoClient(client)

	createUsersCollection(ctx)
	
	createOrganizationsCollection(ctx)
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
func GetOrganizationsCollection() *mongo.Collection {
	return client.Database(dbName).Collection("organizations")
}

func createUsersCollection(ctx context.Context) {
	// Get the MongoDB database
	db := client.Database(dbName)

	// Check if the "users" collection already exists
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": "users"})
	if err != nil {
		log.Fatal(err)
	}

	// If the "users" collection doesn't exist, create it
	if len(collections) == 0 {
		// Specify options for the "users" collection
		collectionOptions := options.CreateCollection().SetValidator(bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"name", "email", "password"}, // Add other required fields if needed
				"properties": bson.M{
					"name": bson.M{
						"bsonType":    "string",
						"description": "must be a string and is required",
					},
					"email": bson.M{
						"bsonType":    "string",
						"description": "must be a string and is required",
					},
					"password": bson.M{
						"bsonType":    "string",
						"description": "must be a string and is required",
					},
					// Add other properties if needed
				},
			},
		})

		// Create the "users" collection with options
		err := db.CreateCollection(ctx, "users", collectionOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Create a unique index on the "email" field
		indexModel := mongo.IndexModel{
			Keys: bson.M{
				"email": 1,
			},
			Options: options.Index().SetUnique(true),
		}

		_, err = db.Collection("users").Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Created 'users' collection with a unique index on 'email' field and required fields.")
	} else {
		log.Println("The 'users' collection already exists.")
	}
}

func createOrganizationsCollection(ctx context.Context) {
	// Get the MongoDB database
	db := client.Database(dbName)

	// Check if the "organizations" collection already exists
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": "organizations"})
	if err != nil {
		log.Fatal(err)
	}

	// If the "organizations" collection doesn't exist, create it
	if len(collections) == 0 {
		// Specify options for the "organizations" collection
		collectionOptions := options.CreateCollection()

		// Create the "organizations" collection with options
		err := db.CreateCollection(ctx, "organizations", collectionOptions)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Created 'organizations' collection.")
	} else {
		log.Println("The 'organizations' collection already exists.")
	}
}


// DropDB drops the MongoDB database.
func DropDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Database(dbName).Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Dropped the database!")
}

// DisconnectClient disconnects the MongoDB client.
func DisconnectClient() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Disconnected the client!")
}
