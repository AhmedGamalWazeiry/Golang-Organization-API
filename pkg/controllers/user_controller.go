package controllers

import (
	"errors"
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

func CreateUser(user models.UserRegister) (int, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	
	user.Password = hashedPassword
	err = repository.CreateUser(user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func AuthenticateUser(userRequest models.UserLoginRequest) (int, *models.User, error) {
	user, err := repository.GetUserByEmail(userRequest.Email)

	if err != nil || user == nil {
		return http.StatusUnauthorized, nil, errors.New("invalid credentials")
	}

	if err := utils.ComparePasswordHash(user.Password, userRequest.Password); err != nil {
		return http.StatusUnauthorized, nil, errors.New("invalid credentials")
	}

	return http.StatusOK, user, nil
}

func GenerateTokenPair(user models.User) (int, string, string, error) {
	accessToken, err := utils.GenerateToken(user, accessTokenExpireMinutes)
	if err != nil {
		return http.StatusInternalServerError, "", "", err
	}
	refreshToken, err := utils.GenerateToken(user, refreshTokenExpireMinutes)
	if err != nil {
		return http.StatusInternalServerError, "", "", err
	}

	return http.StatusOK, accessToken, refreshToken, nil
}

func GenerateTokenPairByRefreshToken(token string) (int, string, string, error) {
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return http.StatusUnauthorized, "", "", errors.New("invalid refresh token")
	}

	user, err := repository.GetUserByID(claims.UserID)
	if err != nil {
		return http.StatusNotFound, "", "", err
	}

	accessToken, err := utils.GenerateToken(*user, accessTokenExpireMinutes)
	if err != nil {
		return http.StatusInternalServerError, "", "", err
	}
	refreshToken, err := utils.GenerateToken(*user, refreshTokenExpireMinutes)
	if err != nil {
		return http.StatusInternalServerError, "", "", err
	}

	return http.StatusOK, accessToken, refreshToken, nil
}
