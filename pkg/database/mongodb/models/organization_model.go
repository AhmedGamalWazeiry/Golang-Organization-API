package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Organization struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty" json:"organization_id"`
	Name                string `bson:"name" binding:"required" json:"name"`
	Description         string `bson:"description" binding:"required" json:"description"`
	OrganizationMembers []Member `bson:"organization_members" json:"organization_members"`
}