package routes

import (
	"github.com/gin-gonic/gin"
	"org.com/org/pkg/api/handlers"
	"org.com/org/pkg/api/middleware"
)

// InitUserRoutes initializes user-related routes.
func InitUserRoutes(router *gin.Engine) {
	userGroup := router.Group("")
	{
		userGroup.POST("/signup", handlers.Register)
		userGroup.POST("/signin", handlers.Login)
		userGroup.POST("/revoke-refresh-token",middleware.Auth(), handlers.RevokeRefreshToken)
		userGroup.POST("/refresh-token", handlers.RefreshToken)
		
	}
}