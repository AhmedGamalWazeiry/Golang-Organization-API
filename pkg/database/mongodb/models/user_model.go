package models

type User struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name" bson:"name" binding:"required"`
	Email    string `json:"email" bson:"email"  binding:"required"`
	Password string `json:"password" bson:"password"  binding:"required"`
}