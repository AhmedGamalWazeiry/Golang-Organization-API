package controllers

import (
	"errors"
	"log"
	"net/http"

	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/database/mongodb/repository"
)

func InsertOrganizationController(orgView models.OrganizationView, userID string) (int, string, error) {
	user,err := repository.GetUserByID(userID)
	if err != nil {
		return http.StatusNotFound, "", errors.New("Can't find the user with this ID")
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
		log.Printf("Error creating organization: %v\n", err)
		return http.StatusInternalServerError, "", err
	}

	return http.StatusCreated, orgID, nil
}

func GetOrganizationByIDController(orgID string, userEmail string) (int, *models.Organization, error) {
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return http.StatusNotFound, nil, errors.New("Can't find the organization with this ID")
	}

	isMember := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail {
			isMember = true
			break
		}
	}

	if !isMember {
		return http.StatusUnauthorized, nil, errors.New("you are not authorized to access this organization's information")
	}

	return http.StatusOK, organization, nil
}

func GetAllUserOrganizationsController(userEmail string) (int, []models.Organization, error) {
	allOrganizations, err := repository.GetAllOrganizations()
	if err != nil {
		return http.StatusNotFound, nil, errors.New("No organizations exist to show.")
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

func UpdateOrganizationController(orgID string, userEmail string, orgView models.OrganizationView) (int, error) {
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return http.StatusNotFound, errors.New("Can't find the organization with this ID")
	}

	isMember := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail {
			isMember = true
			break
		}
	}

	if !isMember {
		return http.StatusUnauthorized, errors.New("you are not authorized to update this organization's information")
	}

	organization.Name = orgView.Name
	organization.Description = orgView.Description

	err = repository.UpdateOrganization(*organization)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func DeleteOrganizationController(orgID string, userEmail string) (int, error) {
	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return http.StatusNotFound, errors.New("Can't find the organization with this ID")
	}

	isMember := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == userEmail {
			isMember = true
			break
		}
	}

	if !isMember {
		return http.StatusUnauthorized, errors.New("you are not authorized to delete this organization's information")
	}

	err = repository.DeleteOrganization(orgID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func InviteUserController(orgID string, userID string, inviteEmail string) (int, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return http.StatusNotFound, err
	}

	organization, err := repository.GetOrganizationByID(orgID)
	if err != nil {
		return http.StatusNotFound, err
	}

	isAdmin := false
	for _, member := range organization.OrganizationMembers {
		if member.Email == user.Email && member.AccessLevel == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return http.StatusUnauthorized, errors.New("you are not authorized to invite users to this organization")
	}

	for _, member := range organization.OrganizationMembers {
		if member.Email == inviteEmail {
			return http.StatusConflict, errors.New("the user is already a member of this organization")
		}
	}

	inviteUser, err := repository.GetUserByEmail(inviteEmail)
	if err != nil {
		return http.StatusNotFound, err
	}

	err = repository.AddUserToOrganization(orgID, *inviteUser)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}