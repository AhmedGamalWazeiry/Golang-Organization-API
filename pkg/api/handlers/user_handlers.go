package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"org.com/org/pkg/controllers"
	"org.com/org/pkg/database/mongodb/models"
	"org.com/org/pkg/utils"
)

// Register handles the POST request for creating a new user.
func Register(c *gin.Context) {
	var user models.UserRegister

	// Bind JSON data to the 'user' variable
	if err := c.ShouldBindJSON(&user); err != nil {
		errorMessage:= utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	emailValidationMessage := utils.ValidateEmail(user.Email)
	if emailValidationMessage != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": emailValidationMessage})
		return
	}
	
	// Create user using controller
	statusCode, err := controllers.CreateUser(user)

	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(statusCode, gin.H{"message": "Congratulations! You've successfully signed up and joined our community."})
}

func Login(c *gin.Context) {
	var user models.UserLoginRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		errorMessage:= utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	emailValidationMessage := utils.IsValidEmail(user.Email)

	if !emailValidationMessage  {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	statusCode, authenticated_user, err := controllers.AuthenticateUser(user)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	statusCode, accessToken, refreshToken, err := controllers.GenerateTokenPair(*authenticated_user)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Login successful",
		"access_token":    accessToken,
		"refresh_token":   refreshToken,
	})
}

func RevokeRefreshToken(c *gin.Context) {
	var RefreshTokenRequest models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&RefreshTokenRequest); err != nil {
		errorMessage:= utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	refreshToken:= RefreshTokenRequest.RefreshToken
	
	err := utils.BlacklistToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to revoke the refresh token; it appears that the refresh token may be invalid."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token revoked successfully"})
}

func RefreshToken(c *gin.Context) {
	var RefreshTokenRequest models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&RefreshTokenRequest); err != nil {
		errorMessage:= utils.ExtractErrorMessage(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	refreshToken:= RefreshTokenRequest.RefreshToken

	ok, err := utils.IsTokenBlacklisted(refreshToken)

	if ok || err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Refresh Token"})
			return
		}
    
	statusCode, accessToken, newRefreshToken, err := controllers.GenerateTokenPairByRefreshToken(refreshToken)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

    err = utils.BlacklistToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create Refresh Token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":         "Token refreshed successfully",
		"access_token":    accessToken,
		"refresh_token":   newRefreshToken,
	})
}

func Test(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message":         "TEST auth successfully",
	
	})
	return
}