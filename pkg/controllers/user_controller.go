package controllers

import (
	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/database/mongodb/repository"
)

// "org.com/org/pkg/database/mongodb/repository"

func CreateUser(user models.User) (string, error) {
	// Implementation using repository function to create a new user
	return repository.CreateUser(user)
}