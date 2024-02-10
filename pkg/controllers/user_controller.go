package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/database/mongodb/repository"

	"org.com/org/pkg/utils"
)

const (
	accessTokenExpireMinutes  = 100
	refreshTokenExpireMinutes = 200
)

var signingKey = []byte("your-secret-key")

// CreateUser creates a new user.
func CreateUser(user models.UserRegister) (int, error) {
	
	isEmailExist, err := repository.IsEmailExists(user.Email)

	if err != nil {
		fmt.Println(err)
		return http.StatusInternalServerError, errors.New("Failed to create user")
	}

	if isEmailExist {
		return http.StatusBadRequest, errors.New("This Email already exists, please use another one")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed to create user")
	}
	
	user.Password = hashedPassword
	err = repository.CreateUser(user)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed to create user")
	}

	return http.StatusCreated, nil
}

// AuthenticateUser authenticates a user.
func AuthenticateUser(userRequest models.UserLoginRequest) (int, *models.User, error) {
	user, err := repository.GetUserByEmail(userRequest.Email)

	if err != nil || user == nil {
		return http.StatusUnauthorized, nil, errors.New("Invalid credentials")
	}

	if err := utils.ComparePasswordHash(user.Password, userRequest.Password); err != nil {
		return http.StatusUnauthorized, nil, errors.New("Invalid credentials")
	}

	return http.StatusOK, user, nil
}

// GenerateTokenPair generates an access token and a refresh token for a user.
func GenerateTokenPair(user models.User) (int, string, string, error) {
	accessToken, err := utils.GenerateToken(user, accessTokenExpireMinutes)
	if err != nil {
		return http.StatusInternalServerError, "", "", errors.New("Failed to generate access token")
	}
	refreshToken, err := utils.GenerateToken(user, refreshTokenExpireMinutes)
	if err != nil {
		return http.StatusInternalServerError, "", "", errors.New("Failed to generate refresh token")
	}

	return http.StatusOK, accessToken, refreshToken, nil
}

// GenerateTokenPairByRefreshToken generates a new pair of tokens using a refresh token.
func GenerateTokenPairByRefreshToken(token string) (int, string, string, error) {
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return http.StatusUnauthorized, "", "", errors.New("Invalid refresh token")
	}

	user, err := repository.GetUserByID(claims.UserID)
	if err != nil {
		return http.StatusNotFound, "", "", errors.New("User not found")
	}

	accessToken, err := utils.GenerateToken(*user, accessTokenExpireMinutes)
	if err != nil {
		return http.StatusInternalServerError, "", "", errors.New("Failed to generate access token")
	}
	refreshToken, err := utils.GenerateToken(*user, refreshTokenExpireMinutes)
	if err != nil {
		return http.StatusInternalServerError, "", "", errors.New("Failed to generate refresh token")
	}

	return http.StatusOK, accessToken, refreshToken, nil
}
