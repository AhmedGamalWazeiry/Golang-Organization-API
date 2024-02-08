package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"org.com/org/pkg/utils"
)

// Auth is a middleware function that performs user authorization based on the "Authorization" header.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the "Authorization" header
		authHeader := c.GetHeader("Authorization")

		// Check if the header is missing
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Check if the header has the correct format (Bearer <token>)
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		// Extract the access token
		accessToken := headerParts[1]

		// Verify the token
		claims, err := utils.VerifyToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set user ID from claims for further use in the handlers
		c.Set("user_id", claims.UserID)

		// Continue with the next middleware/handler
		c.Next()
	}
}
