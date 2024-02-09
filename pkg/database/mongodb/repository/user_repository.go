package repository

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"org.com/org/pkg/database/mongodb"
	"org.com/org/pkg/database/mongodb/models"
)

// CreateUser creates a new user in the database.
func CreateUser(user models.UserRegister) error {
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

func IsEmailExists(email string) (bool, error) {
    collection := mongodb.GetUsersCollection()

    filter := bson.M{"email": email}
    count, err := collection.CountDocuments(context.Background(), filter)
    if err != nil {
        return false, err
    }

    return count > 0, nil
}

// GetUserByID retrieves a user by ID from the database.
func GetUserByID(userID string) (*models.User, error) {
    collection := mongodb.GetUsersCollection()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

    filter := bson.M{"_id": objectID}
    var user models.User
    err = collection.FindOne(context.Background(), filter).Decode(&user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}