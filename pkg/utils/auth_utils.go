package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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
	jti := GenerateUUID()

	// Create token claims
	claims := &models.TokenClaims{
		UserID: user.ID,
		JTI: jti,
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
	
	if err != nil {
		return nil, errors.New("token validation failed")
	}
    
	claims, ok := token.Claims.(*models.TokenClaims)

	if !ok || !token.Valid {
		return nil, errors.New("token validation failed")
	}

	return claims, nil
}

func ExtractToken(c *gin.Context) (string, bool) {
	// Get the "Authorization" header
	authHeader := c.GetHeader("Authorization")

	// Check if the header is missing
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return "", false
	}

	// Check if the header has the correct format (Bearer <token>)
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return "", false
	}

	// Extract the access token
	accessToken := headerParts[1]

	return accessToken, true
}

// BlacklistToken adds the given token to the blacklist with a specific expiration time.
func  BlacklistToken(token string) error {
	claims, err := VerifyToken(token)
	if err != nil {
		return err
	}

	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", claims.JTI)
	keys, err := RedisClient.Keys(context.Background(), "*").Result()
	if err != nil {
		fmt.Println("Error:", err)
		
	}

	// Print the keys
	fmt.Println("Keys:")
	for _, key := range keys {
		fmt.Println(key)
	}

	remainingTime := time.Until(time.Unix(claims.ExpiresAt, 0))

	return RedisClient.Set(ctx, key, true, remainingTime).Err()
}
// IsTokenBlacklisted checks if the given token is blacklisted.
func IsTokenBlacklisted(token string) (bool, error) {
	claims, err := VerifyToken(token)
	if err != nil {
		return false,err
	}

	ctx := context.Background()
	key := fmt.Sprintf("blacklist:%s", claims.JTI)

	result, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return result > 0, nil
}