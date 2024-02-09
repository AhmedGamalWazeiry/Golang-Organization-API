package models

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name" binding:"required"`
	Email    string `json:"email" bson:"email"  binding:"required"`
	Password string `json:"password" bson:"password"  binding:"required"`
}

// LoginRequest represents the model for the login request.
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
