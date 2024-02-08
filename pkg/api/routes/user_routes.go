package routes

import (
	"github.com/gin-gonic/gin"
	"org.com/org/pkg/api/handlers"
	"org.com/org/pkg/api/middleware"
)

// InitUserRoutes initializes user-related routes.
func InitUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/signup", handlers.Register)
		userGroup.POST("/signin", handlers.Login)
		userGroup.POST("/signout", handlers.Logout)
		userGroup.GET("/test",middleware.Auth(), handlers.Logout)
	}
}