package routes

import (
	"github.com/gin-gonic/gin"
	"org.com/org/pkg/api/handlers"
)

// InitUserRoutes initializes user-related routes.
func InitUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("/:id", handlers.GetUserByID)
		userGroup.POST("/", handlers.CreateUser)
		userGroup.PUT("/:id", handlers.UpdateUser)
		userGroup.DELETE("/:id", handlers.DeleteUser)
	}
}