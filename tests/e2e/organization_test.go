package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"org.com/org/pkg/database/mongodb/models"
)

var (
	organizationID string
)

func TestCreateOrganizationHandler(t *testing.T) {
	// Create a new organization request for testing
	organizationRequest := models.Organization{
		Name:        "Test Organization",
		Description: "This is a test organization",
	}

	// Convert the organization request to JSON
	organizationJSON, _ := json.Marshal(organizationRequest)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/organization", bytes.NewBuffer(organizationJSON))
	req.Header.Set("Content-Type", "application/json")

	// Add the access token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)

	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	
	organizationID = response["organization_id"]
	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Assert that the response body contains the expected content
	assert.Contains(t, resp.Body.String(), "organization_id")
}

func TestGetOrganizationByIDHandler(t *testing.T) {
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/organization/"+organizationID, nil)
	req.Header.Set("Content-Type", "application/json")

	// Add the access token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert that the response body contains the expected content
	assert.Contains(t, resp.Body.String(), "organization_id")
}

func TestGetAllUserOrganizationsHandler(t *testing.T) {
	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/organization", nil)
	req.Header.Set("Content-Type", "application/json")

	// Add the access token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert that the response body contains the expected content
	var organizations []models.Organization
	err := json.Unmarshal(resp.Body.Bytes(), &organizations)
	assert.Nil(t, err)
	assert.NotEmpty(t, organizations)
}

func TestUpdateOrganizationHandler(t *testing.T) {
	// Create a new organization request for testing
	organizationRequest := models.Organization{
		Name: "Updated Test Organization",
		Description: "This is an updated test organization",
	}

	// Convert the organization request to JSON
	organizationJSON, _ := json.Marshal(organizationRequest)

	// Create a new HTTP request
	req, _ := http.NewRequest("PUT", "/organization/"+organizationID, bytes.NewBuffer(organizationJSON))
	req.Header.Set("Content-Type", "application/json")

	// Add the access token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert that the response body contains the expected content
	var organization models.Organization
	err := json.Unmarshal(resp.Body.Bytes(), &organization)
	assert.Nil(t, err)
	assert.Equal(t, "Updated Test Organization", organization.Name)
	assert.Equal(t, "This is an updated test organization", organization.Description)
}

func TestInviteUserHandler(t *testing.T) {
	// Create a new invite request for testing
	inviteRequest := models.Invite{
		Email: "testuserInvite@example.com",
	}

	// Convert the invite request to JSON
	inviteJSON, _ := json.Marshal(inviteRequest)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/organization/"+organizationID+"/invite", bytes.NewBuffer(inviteJSON))
	req.Header.Set("Content-Type", "application/json")

	// Add the access token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert that the response body contains the expected content
	assert.Contains(t, resp.Body.String(), "User invited successfully")
}

func TestDeleteOrganizationHandler(t *testing.T) {
	// Create a new HTTP request
	req, _ := http.NewRequest("DELETE", "/organization/"+organizationID, nil)
	req.Header.Set("Content-Type", "application/json")

	// Add the access token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert that the response body contains the expected content
	assert.Contains(t, resp.Body.String(), "Organization deleted successfully")
}
