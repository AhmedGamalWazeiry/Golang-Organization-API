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

func TestUserRegistration(t *testing.T) {
	// Create a new user for testing
	user := models.User{
		Name:     "Test User",
		Email:    "testuser1@example.com",
		Password: "password1234",
	}

	// Convert the user to JSON
	userJSON, _ := json.Marshal(user)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Assert that the response body contains the expected content
	assert.Contains(t, resp.Body.String(), "Congratulations! You've successfully signed up and joined our community.")
}

func TestUserLogin(t *testing.T) {
	

	// Create a new user for testing
	user := models.UserLoginRequest{
		Email:    "testuser1@example.com",
		Password: "password1234",
	}

	// Convert the user to JSON
	userJSON, _ := json.Marshal(user)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert that the response body contains the expected content
	assert.Contains(t, resp.Body.String(), "Login successful")
}



func TestRefreshToken(t *testing.T) {
	// Create a new refresh token request for testing
	refreshTokenRequest := models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	// Convert the refresh token request to JSON
	refreshTokenJSON, _ := json.Marshal(refreshTokenRequest)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/refresh-token", bytes.NewBuffer(refreshTokenJSON))
	req.Header.Set("Content-Type", "application/json")

	// Record the HTTP response
	resp := httptest.NewRecorder()
    
	
	// Serve the HTTP request
	router.ServeHTTP(resp, req)
     
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	
	accessToken = response["access_token"]
	refreshToken = response["refresh_token"]
   
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert that the response body contains the expected content
	assert.Contains(t, resp.Body.String(), "Token refreshed successfully")
}

func TestRevokeRefreshToken(t *testing.T) {
	// Create a new refresh token request for testing
	refreshTokenRequest := models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	// Convert the refresh token request to JSON
	refreshTokenJSON, _ := json.Marshal(refreshTokenRequest)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/revoke-refresh-token", bytes.NewBuffer(refreshTokenJSON))
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
	assert.Contains(t, resp.Body.String(), "Refresh token revoked successfully")
}



