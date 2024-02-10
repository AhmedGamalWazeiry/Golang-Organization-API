package controllers

import (
	"errors"
	"net/http"

	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/database/mongodb/repository"
)

// InsertOrganizationController creates a new organization.
func InsertOrganizationController(orgView models.OrganizationView, userID string) (int, string, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return http.StatusNotFound, "", errors.New("User not found")
	}

	member := models.Member{
		Name:       user.Name,
		Email:      user.Email,
		AccessLevel: "admin",
	}

	organization := models.Organization{
		Name:                 orgView.Name,
		Description:          orgView.Description,
		OrganizationMembers: []models.Member{member},
	}

	orgID, err := repository.CreateOrganization(organization)
	if err != nil {
		return http.StatusInternalServerError, "", errors.New("Failed to create organization")
	}

	return http.StatusCreated, orgID, nil
}

// GetOrganizationByIDController retrieves an organization by its ID.
func GetOrganizationByIDController(orgID string, userEmail string) (int, *models.Organization, error) {
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return http.StatusNotFound, nil, errors.New("Organization not found")
	}

	isMember := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail {
			isMember = true
			break
		}
	}

	if !isMember {
		return http.StatusForbidden, nil, errors.New("You don't have access to perform this action.")
	}

	return http.StatusOK, organization, nil
}

// GetAllUserOrganizationsController retrieves all organizations a user is part of.
func GetAllUserOrganizationsController(userEmail string) (int, []models.Organization, error) {
	allOrganizations, err := repository.GetAllOrganizations()
	if err != nil {
		return http.StatusNotFound, nil, errors.New("No organizations found")
	}

	var userOrganizations []models.Organization
	for _, org := range allOrganizations {
		for _, member := range org.OrganizationMembers {
			if member.Email == userEmail {
				userOrganizations = append(userOrganizations, org)
				break
			}
		}
	}

	return http.StatusOK, userOrganizations, nil
}

// UpdateOrganizationController updates an organization's information.
func UpdateOrganizationController(orgID string, userEmail string, orgView models.OrganizationView) (int, error) {
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return http.StatusNotFound, errors.New("Organization not found")
	}

	isAdmin  := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail && member.AccessLevel == "admin" {
			isAdmin  = true
			break
		}
	}

	if !isAdmin  {
		return http.StatusForbidden, errors.New("You don't have access to perform this action.")
	}

	organization.Name = orgView.Name
	organization.Description = orgView.Description

	err = repository.UpdateOrganization(*organization)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed to update organization")
	}

	return http.StatusOK, nil
}

// DeleteOrganizationController deletes an organization.
func DeleteOrganizationController(orgID string, userEmail string) (int, error) {
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return http.StatusNotFound, errors.New("Organization not found")
	}

	isAdmin  := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail && member.AccessLevel == "admin" {
			isAdmin  = true
			break
		}
	}

	if !isAdmin  {
		return http.StatusForbidden, errors.New("You don't have access to perform this action.")
	}

	err = repository.DeleteOrganization(orgID)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed to delete organization")
	}

	return http.StatusOK, nil
}

// InviteUserController invites a user to an organization.
func InviteUserController(orgID string, userID string, inviteEmail string) (int, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return http.StatusNotFound, errors.New("User not found")
	}

	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return http.StatusNotFound, errors.New("Organization not found")
	}

	isAdmin := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == user.Email && member.AccessLevel == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return http.StatusForbidden, errors.New("You don't have access to perform this action.")
	}

	for _, member := range organization.OrganizationMembers {
		if member.Email == inviteEmail {
			return http.StatusConflict, errors.New("User already a member of this organization")
		}
	}

	inviteUser, err := repository.GetUserByEmail(inviteEmail)
	if err != nil {
		return http.StatusNotFound, errors.New("User not found")
	}

	err = repository.AddUserToOrganization(orgID, *inviteUser)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Failed to add user to organization")
	}

	return http.StatusOK, nil
}
