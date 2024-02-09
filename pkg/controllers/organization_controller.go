package controllers

import (
	"errors"
	"log"

	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/database/mongodb/repository"
)

func InsertOrganizationController(orgView models.OrganizationView, userID string) (string, error) {
	// Get the user collection
	

	// Convert the userID to an ObjectID
	user,err := repository.GetUserByID(userID)
	if err != nil {
		return "", errors.New("Can't find the user with this ID")
	}

	// Create a new member with the user data and admin access level
	member := models.Member{
		Name:       user.Name,
		Email:      user.Email,
		AccessLevel: "admin",
	}

	// Create a new organization with the OrganizationView and added member
	organization := models.Organization{
		Name:                 orgView.Name,
		Description:          orgView.Description,
		OrganizationMembers: []models.Member{member},
	}

	// Use the existing CreateOrganization function to insert the new organization
	orgID, err := repository.CreateOrganization(organization)
	if err != nil {
		log.Printf("Error creating organization: %v\n", err)
		return "", err
	}

	return orgID, nil
}

func GetOrganizationByIDController(orgID string, userEmail string) (*models.Organization, error) {
	

	// Get the organization object from GetOrganizationByID in the repo
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return nil, errors.New("Can't find the organization with this ID")
	}

	// Make sure that the user email exists in the members of the organization
	isMember := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail {
			isMember = true
			break
		}
	}

	// If the user is not a member in the organization, return an error
	if !isMember {
		return nil, errors.New("you are not authorized to access this organization's information")
	}

	return organization, nil
}

func GetAllUserOrganizationsController(userEmail string) ([]models.Organization, error) {
	// Get all organizations from the repo
	allOrganizations, err := repository.GetAllOrganizations()
	if err != nil {
		return nil,  errors.New("No organizations exist to show.")
	}

	// Filter the organizations to include only those where the user is a member
	var userOrganizations []models.Organization
	for _, org := range allOrganizations {
		for _, member := range org.OrganizationMembers {
			if member.Email == userEmail {
				userOrganizations = append(userOrganizations, org)
				break
			}
		}
	}

	return userOrganizations, nil
}

func UpdateOrganizationController(orgID string, userEmail string, orgView models.OrganizationView) error {
	// Get the user by the userID in the repo
	

	// Get the organization object from GetOrganizationByID in the repo
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return errors.New("Can't find the organization with this ID")
	}

	// Make sure that the user email exists in the members of the organization
	isMember := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail {
			isMember = true
			break
		}
	}

	// If the user is not a member in the organization, return an error
	if !isMember {
		return errors.New("you are not authorized to update this organization's information")
	}

	// Update the organization with the new data
	organization.Name = orgView.Name
	organization.Description = orgView.Description

	// Use the existing UpdateOrganization function to update the organization
	err = repository.UpdateOrganization(*organization)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOrganizationController(orgID string, userEmail string) error {

	// Get the organization object from GetOrganizationByID in the repo
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return  errors.New("Can't find the organization with this ID")
	}

	// Make sure that the user email exists in the members of the organization
	isMember := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail {
			isMember = true
			break
		}
	}

	// If the user is not a member in the organization, return an error
	if !isMember {
		return errors.New("you are not authorized to delete this organization's information")
	}

	// Use the existing DeleteOrganization function to delete the organization
	err = repository.DeleteOrganization(orgID)
	if err != nil {
		return err
	}

	return nil
}

func InviteUserController(orgID string, userID string, inviteEmail string) error {
	// Get the user by the userID in the repo
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Get the organization object from GetOrganizationByID in the repo
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return err
	}

	// Make sure that the user email exists in the members of the organization and the access level is admin
	isAdmin := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == user.Email && member.AccessLevel == "admin" {
			isAdmin = true
			break
		}
	}

	// If the user is not an admin in the organization, return an error
	if !isAdmin {
		return errors.New("you are not authorized to invite users to this organization")
	}

	// Check that the user to invite does not exist in the same organization
	for _, member := range organization.OrganizationMembers {
		if member.Email == inviteEmail {
			return errors.New("the user is already a member of this organization")
		}
	}

	// Get the invited user by the email in the repo
	inviteUser, err := repository.GetUserByEmail(inviteEmail)
	if err != nil {
		return err
	}

	// Add the invited user to the organization
	err = repository.AddUserToOrganization(orgID, *inviteUser)
	if err != nil {
		return err
	}

	return nil
}
