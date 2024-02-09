package repository

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"org.com/org/pkg/database/mongodb"
	"org.com/org/pkg/database/mongodb/models" // replace with the actual package name
)

func CreateOrganization(organization models.Organization) (string, error) {
    collection := mongodb.GetOrganizationsCollection()

    // Inserting a new organization document into the collection
    result, err := collection.InsertOne(context.TODO(), organization)
    if err != nil {
        log.Printf("Error creating organization: %v\n", err)
        return "", err
    }

    // Retrieve the generated ID from the InsertOneResult
    insertedID, ok := result.InsertedID.(primitive.ObjectID)
    if !ok {
        log.Printf("Error converting InsertedID to ObjectID")
        return "", errors.New("error converting InsertedID to ObjectID")
    }

    // Convert ObjectID to string
    idString := insertedID.Hex()

    return idString, nil
}



func GetOrganizationByID(organizationID string) (*models.Organization, error) {
    collection := mongodb.GetOrganizationsCollection()

	objectID, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		return nil, errors.New("invalid organization ID format")
	}

    filter := bson.M{"_id": objectID}
    var organization models.Organization
    err = collection.FindOne(context.Background(), filter).Decode(&organization)
    if err != nil {
        return nil, err
    }

    return &organization, nil
}

// GetAllOrganizations retrieves all organizations from the MongoDB collection.
func GetAllOrganizations() ([]models.Organization, error) {
    collection := mongodb.GetOrganizationsCollection()

    // Define an empty filter to get all documents in the collection
    filter := bson.M{}

    // Find all organizations matching the filter
    cursor, err := collection.Find(context.Background(), filter)
    if err != nil {
        log.Printf("Error retrieving organizations: %v\n", err)
        return nil, err
    }
    defer cursor.Close(context.Background())

    // Iterate through the cursor and decode documents into an array of organizations
    var organizations []models.Organization
    for cursor.Next(context.Background()) {
        var organization models.Organization
        if err := cursor.Decode(&organization); err != nil {
            log.Printf("Error decoding organization: %v\n", err)
            return nil, err
        }
        organizations = append(organizations, organization)
    }

    if err := cursor.Err(); err != nil {
        log.Printf("Error iterating over organizations: %v\n", err)
        return nil, err
    }

    return organizations, nil
}

func UpdateOrganization(organization models.Organization) error {
    collection := mongodb.GetOrganizationsCollection()

    // Create a filter for the organization to update
    filter := bson.M{"_id": organization.ID}

    // Create an update that sets the new organization data
    update := bson.M{
        "$set": bson.M{
            "name": organization.Name,
            "description": organization.Description,
            "organization_members": organization.OrganizationMembers,
        },
    }

    // Update the organization in the collection
    _, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }

    return nil
}

func DeleteOrganization(orgID string) error {
    collection := mongodb.GetOrganizationsCollection()

    // Convert the organizationID to an ObjectID
    objectID, err := primitive.ObjectIDFromHex(orgID)
    if err != nil {
        return errors.New("invalid organization ID format")
    }

    // Delete the organization from the collection
    _, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
    if err != nil {
        return err
    }

    return nil
}

func AddUserToOrganization(orgID string, user models.User) error {
    collection := mongodb.GetOrganizationsCollection()

    // Convert the organizationID to an ObjectID
    objectID, err := primitive.ObjectIDFromHex(orgID)
    if err != nil {
        return errors.New("invalid organization ID format")
    }

    // Create a new member with the user name, email, and default access level
    newMember := models.Member{
        Name: user.Name,
        Email: user.Email,
        AccessLevel: "member",  // Set the default access level for new members
    }

    // Add the new member to the organization in the collection
    _, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$push": bson.M{"organization_members": newMember}})
    if err != nil {
        return err
    }

    return nil
}