package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"org.com/org/pkg/database/mongodb/models"
)


var signingKey = []byte("your-secret-key")

func HashPassword(password string) (string, error) {
	// Generate a hashed version of the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	
	return string(hashedPassword), nil
}

func ComparePasswordHash(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
// GenerateUUID generates a new UUID.
func GenerateUUID() string {
	return uuid.New().String()
}


// ParseToken parses and validates the given token.
func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{},func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
func GenerateToken(user models.User,tokenExpireMinutes int) (string, error) {

	expirationTime := time.Now().Add(time.Minute * time.Duration(tokenExpireMinutes)).Unix()

	// Create token claims
	claims := &models.TokenClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Email,
			ExpiresAt: expirationTime,
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


// VerifyToken verifies the access token and returns the claims.
func VerifyToken(tokenString string) (*models.TokenClaims, error) {
	token, err := ParseToken(tokenString)
	fmt.Println(tokenString)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
    
	claims, ok := token.Claims.(*models.TokenClaims)
	fmt.Println(claims)
	fmt.Println(ok)
	if !ok || !token.Valid {
		return nil, errors.New("token validation failed")
	}

	return claims, nil
}