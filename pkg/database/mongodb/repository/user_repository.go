package repository

import (
	// "go.mongodb.org/mongo-driver/mongo"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"org.com/org/pkg/database/mongodb"

	"org.com/org/pkg/database/mongodb/models"
)

// "go.mongodb.org/mongo-driver/bson"

// // GetUserByID retrieves a user by ID from the database.
// func GetUserByID(id string) (models.User, error) {

// 	// Implementation using MongoDB driver to find a user by ID
// }

// CreateUser creates a new user in the database.
func CreateUser(user models.User) (string, error) {
	// Implementation using MongoDB driver to insert a new user
	collection := mongodb.GetUsersCollection()
	// Inserting a new user document into the collection
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return "", err
	}
	// Extracting the created document ID and returning it
	createdID := result.InsertedID.(primitive.ObjectID).Hex()
	return createdID, nil
}

// // UpdateUser updates a user by ID in the database.
// func UpdateUser(id string, updatedUser models.User) error {
// 	// Implementation using MongoDB driver to update a user
// }

// // DeleteUser deletes a user by ID from the database.
// func DeleteUser(id string) error {
// 	// Implementation using MongoDB driver to delete a user
// }