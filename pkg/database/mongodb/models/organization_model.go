package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty"`
	Name                 string `bson:"name" binding:"required"`
	Description          string `bson:"description" binding:"required"`
	OrganizationMembers []Member `bson:"organization_members"`
}

// User represents the structure of the user document in MongoDB
type Member struct {
	Name       string `bson:"name"`
	Email      string `bson:"email"`
	AccessLevel string `bson:"access_level"`
}

type OrganizationView struct {
	Name                 string `bson:"name" binding:"required"`
	Description          string `bson:"description" binding:"required"`
}