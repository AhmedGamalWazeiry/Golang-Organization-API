package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"org.com/org/pkg/database/mongodb"
	"org.com/org/pkg/database/mongodb/models"
)

// CreateUser creates a new user in the database.
func CreateUser(user models.User) error {
	collection := mongodb.GetUsersCollection()

	// Inserting a new user document into the collection
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return err
	}

	return nil
}

// GetUserByEmail retrieves a user by email from the database.
func GetUserByEmail(email string) (*models.User, error) {
	collection := mongodb.GetUsersCollection()

	filter := bson.M{"email": email}
	var user models.User
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
