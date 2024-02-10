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
	statusCode, errorMessage := utils.BindAndValidate(c, &user)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{"error": errorMessage})
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

// Login handles the POST request for user login.
func Login(c *gin.Context) {
	var user models.UserLoginRequest

	statusCode, errorMessage := utils.BindAndValidate(c, &user)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	statusCode, authenticated_user, err := controllers.AuthenticateUser(user)
	if err != nil {
		c.JSON(statusCode, err.Error())
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

// RevokeRefreshToken handles the POST request for revoking a refresh token.
func RevokeRefreshToken(c *gin.Context) {
	var RefreshTokenRequest models.RefreshTokenRequest

	statusCode, errorMessage := utils.BindAndValidate(c, &RefreshTokenRequest)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{"error": errorMessage})
		return
	}

	refreshToken:= RefreshTokenRequest.RefreshToken
	
	err := utils.BlacklistToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to revoke the refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Refresh token revoked successfully"})
}

// RefreshToken handles the POST request for refreshing a token.
func RefreshToken(c *gin.Context) {
	var RefreshTokenRequest models.RefreshTokenRequest

	statusCode, errorMessage := utils.BindAndValidate(c, &RefreshTokenRequest)
	if statusCode != http.StatusOK {
		c.JSON(statusCode, gin.H{"error": errorMessage})
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
		c.JSON(statusCode, gin.H{"error": "Failed to refresh tokens"})
		return
	}

    err = utils.BlacklistToken(refreshToken)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to blacklist old refresh token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":         "Token refreshed successfully",
		"access_token":    accessToken,
		"refresh_token":   newRefreshToken,
	})
}