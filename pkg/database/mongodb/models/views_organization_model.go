package models

type Member struct {
	Name        string `bson:"name" json:"name"`
	Email       string `bson:"email" json:"email"`
	AccessLevel string `bson:"access_level" json:"access_level"`
}

type OrganizationView struct {
	Name        string `bson:"name" binding:"required" json:"name"`
	Description string `bson:"description" binding:"required" json:"description"`
}

type Invite struct {
	Email string `bson:"user_email" binding:"required" json:"user_email"`
}
