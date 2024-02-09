package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"org.com/org/pkg/utils"
)

// Auth is a middleware function that performs user authorization based on the "Authorization" header.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token
		token, ok := utils.ExtractToken(c)
		if !ok {
			c.Abort()
			return
		}

		// Verify the token
		claims, err := utils.VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		ok, err = utils.IsTokenBlacklisted(token)

		if ok || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Set user ID from claims for further use in the handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Subject)

		// Continue with the next middleware/handler
		c.Next()
	}
}
