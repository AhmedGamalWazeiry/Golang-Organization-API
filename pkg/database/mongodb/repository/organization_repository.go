package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"org.com/org/pkg/database/mongodb"
	"org.com/org/pkg/database/mongodb/models"
)

// CreateOrganization inserts a new organization document into the collection.
func CreateOrganization(organization models.Organization) (string, error) {
    collection := mongodb.GetOrganizationsCollection()

    result, err := collection.InsertOne(context.TODO(), organization)
    if err != nil {
        return "", err
    }

    insertedID, ok := result.InsertedID.(primitive.ObjectID)
    if !ok {
        return "", errors.New("Error processing request")
    }

    return insertedID.Hex(), nil
}

// GetOrganizationByID retrieves an organization document by its ID.
func GetOrganizationByID(organizationID string) (*models.Organization, error) {
    collection := mongodb.GetOrganizationsCollection()

    objectID, err := primitive.ObjectIDFromHex(organizationID)
    if err != nil {
        return nil, errors.New("Invalid organization ID format")
    }

    filter := bson.M{"_id": objectID}

    var organization models.Organization
    err = collection.FindOne(context.Background(), filter).Decode(&organization)
    if err != nil {
        return nil, err
    }

    return &organization, nil
}

// GetAllOrganizations retrieves all organization documents from the collection.
func GetAllOrganizations() ([]models.Organization, error) {
    collection := mongodb.GetOrganizationsCollection()

    filter := bson.M{}

    cursor, err := collection.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    var organizations []models.Organization
    for cursor.Next(context.Background()) {
        var organization models.Organization
        if err := cursor.Decode(&organization); err != nil {
            return nil, err
        }
        organizations = append(organizations, organization)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return organizations, nil
}

// UpdateOrganization updates an organization document in the collection.
func UpdateOrganization(organization models.Organization) error {
    collection := mongodb.GetOrganizationsCollection()

    filter := bson.M{"_id": organization.ID}

    update := bson.M{
        "$set": bson.M{
            "name": organization.Name,
            "description": organization.Description,
            "organization_members": organization.OrganizationMembers,
        },
    }

    _, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }

    return nil
}

// DeleteOrganization deletes an organization document from the collection.
func DeleteOrganization(orgID string) error {
    collection := mongodb.GetOrganizationsCollection()

    objectID, err := primitive.ObjectIDFromHex(orgID)
    if err != nil {
        return errors.New("Invalid organization ID format")
    }

    _, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
    if err != nil {
        return err
    }

    
    return nil
}


// AddUserToOrganization adds a user to an organization as a member.
func AddUserToOrganization(orgID string, user models.User) error {
    collection := mongodb.GetOrganizationsCollection()

    objectID, err := primitive.ObjectIDFromHex(orgID)
    if err != nil {
        return errors.New("The organization ID you provided doesn't seem to be correct. Please check and try again.")
    }

    newMember := models.Member{
        Name: user.Name,
        Email: user.Email,
        AccessLevel: "member", 
    }

    update := bson.M{"$push": bson.M{"organization_members": newMember}}

    _, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
    if err != nil {
        return err
    }

    return nil
}
