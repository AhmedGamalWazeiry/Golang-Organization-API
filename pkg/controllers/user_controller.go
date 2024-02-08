package controllers

import (
	"errors"

	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/database/mongodb/repository"

	"org.com/org/pkg/utils"
)
const (
	accessTokenExpireMinutes  = 5
	refreshTokenExpireMinutes = 20
)

var signingKey = []byte("your-secret-key")

func CreateUser(user models.User) error {
	// Implementation using repository function to create a new user

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	
	// Set the hashed password in the user object
	user.Password = hashedPassword
	return repository.CreateUser(user)
}

// AuthenticateUser handles the authentication logic for a user.
func AuthenticateUser(userRequest models.UserLoginRequest) (*models.User, error) {
	// Check if the email exists
	user, err := repository.GetUserByEmail(userRequest.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare the hashed password with the provided password
	if err := utils.ComparePasswordHash(user.Password, userRequest.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GenerateTokenPair generates a new access and refresh token pair.
func GenerateTokenPair(user models.User) (string,string, error) {
	
	accessToken,err := utils.GenerateToken(user,accessTokenExpireMinutes)
	if err != nil {
		return "","", err
	}
	refreshToken,err := utils.GenerateToken(user,refreshTokenExpireMinutes)
	if err != nil {
		return "","", err
	}
	

	return accessToken,refreshToken, nil
}



