package models

import "github.com/dgrijalva/jwt-go"

// TokenClaims represents the claims for JWT tokens.
type TokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
	JTI    string `json:"jti"`
}