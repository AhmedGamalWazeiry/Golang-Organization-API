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
    // Get the organizations collection from the database.
    collection := mongodb.GetOrganizationsCollection()

    // Insert the organization into the collection.
    result, err := collection.InsertOne(context.TODO(), organization)
    if err != nil {
        return "", err
    }

    // Retrieve the generated ID from the InsertOneResult.
    insertedID, ok := result.InsertedID.(primitive.ObjectID)
    if !ok {
        return "", errors.New("Error processing request")
    }

    // Convert ObjectID to string and return.
    return insertedID.Hex(), nil
}

// GetOrganizationByID retrieves an organization document by its ID.
func GetOrganizationByID(organizationID string) (*models.Organization, error) {
    // Get the organizations collection from the database.
    collection := mongodb.GetOrganizationsCollection()

    // Convert the organizationID to an ObjectID.
    objectID, err := primitive.ObjectIDFromHex(organizationID)
    if err != nil {
        return nil, errors.New("Invalid organization ID format")
    }

    // Define a filter to find the organization by its ID.
    filter := bson.M{"_id": objectID}

    // Find the organization in the collection.
    var organization models.Organization
    err = collection.FindOne(context.Background(), filter).Decode(&organization)
    if err != nil {
        return nil, err
    }

    // Return the found organization.
    return &organization, nil
}

// GetAllOrganizations retrieves all organization documents from the collection.
func GetAllOrganizations() ([]models.Organization, error) {
    // Get the organizations collection from the database.
    collection := mongodb.GetOrganizationsCollection()

    // Define an empty filter to get all documents in the collection.
    filter := bson.M{}

    // Find all organizations in the collection.
    cursor, err := collection.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())

    // Iterate through the cursor and decode documents into an array of organizations.
    var organizations []models.Organization
    for cursor.Next(context.Background()) {
        var organization models.Organization
        if err := cursor.Decode(&organization); err != nil {
            return nil, err
        }
        organizations = append(organizations, organization)
    }

    // Check for errors during cursor iteration.
    if err := cursor.Err(); err != nil {
        return nil, err
    }

    // Return the array of organizations.
    return organizations, nil
}

// UpdateOrganization updates an organization document in the collection.
func UpdateOrganization(organization models.Organization) error {
    // Get the organizations collection from the database.
    collection := mongodb.GetOrganizationsCollection()

    // Define a filter to find the organization by its ID.
    filter := bson.M{"_id": organization.ID}

    // Define an update that sets the new organization data.
    update := bson.M{
        "$set": bson.M{
            "name": organization.Name,
            "description": organization.Description,
            "organization_members": organization.OrganizationMembers,
        },
    }

    // Update the organization in the collection.
    _, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }

    return nil
}

// DeleteOrganization deletes an organization document from the collection.
func DeleteOrganization(orgID string) error {
    // Get the organizations collection from the database.
    collection := mongodb.GetOrganizationsCollection()

    // Convert the organizationID to an ObjectID.
    objectID, err := primitive.ObjectIDFromHex(orgID)
    if err != nil {
        return errors.New("Invalid organization ID format")
    }

    // Delete the organization from the collection.
    _, err = collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
    if err != nil {
        return err
    }

    
    return nil
}


// AddUserToOrganization adds a user to an organization as a member.
func AddUserToOrganization(orgID string, user models.User) error {
    // Get the organizations collection from the database.
    collection := mongodb.GetOrganizationsCollection()

    // Convert the organizationID to an ObjectID.
    objectID, err := primitive.ObjectIDFromHex(orgID)
    if err != nil {
        return errors.New("The organization ID you provided doesn't seem to be correct. Please check and try again.")
    }

    // Create a new member with the user's name, email, and a default access level.
    newMember := models.Member{
        Name: user.Name,
        Email: user.Email,
        AccessLevel: "member", 
    }

    // Define an update operation to add the new member to the organization.
    update := bson.M{"$push": bson.M{"organization_members": newMember}}

    // Update the organization in the collection by adding the new member.
    _, err = collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
    if err != nil {
        return err
    }

    return nil
}
