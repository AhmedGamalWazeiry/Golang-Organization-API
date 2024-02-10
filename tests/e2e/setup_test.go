package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"org.com/org/pkg/api/routes"
	"org.com/org/pkg/database/mongodb"
	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/utils"
)

var (
	router       *gin.Engine
	accessToken  string
	refreshToken string
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()

	os.Exit(code)
}
func setup() {
	
	dbName := fmt.Sprintf("test_db_%d", time.Now().UnixNano())

	mongodb.InitDB(dbName,"mongodb://localhost:27017")

	utils.InitRedis("localhost:6379","",0)
	
	router = gin.Default()
	routes.InitUserRoutes(router)
	routes.InitOrganizationRoutes(router)

	createTestUser("testuser@example.com")
	createTestUser("testuserInvite@example.com")
	loginTestUser()
}

func teardown() {
 mongodb.DropDB()
 mongodb.DisconnectClient()
}


func createTestUser(email string) {
	testUser := models.User{
		Name:     "Test User",
		Email:    email,
		Password: "password123",
	}
	// Convert the user to JSON
	userJSON, _ := json.Marshal(testUser)

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Record the HTTP response
	resp := httptest.NewRecorder()

	// Serve the HTTP request
	router.ServeHTTP(resp, req)
}

func loginTestUser() {
	// Login the test user
	userLoginRequest := models.UserLoginRequest{
		Email:    "testuser@example.com",
		Password: "password123",
	}
	userLoginJSON, _ := json.Marshal(userLoginRequest)
	req, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(userLoginJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Parse the response and store the tokens
	var response map[string]string
	json.Unmarshal(resp.Body.Bytes(), &response)
	accessToken = response["access_token"]
	refreshToken = response["refresh_token"]
}