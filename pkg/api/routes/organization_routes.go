package routes

import (
	"github.com/gin-gonic/gin"
	"org.com/org/pkg/api/handlers"
	"org.com/org/pkg/api/middleware"
)

// InitUserRoutes initializes user-related routes.
func InitOrganizationRoutes(router *gin.Engine) {
	userGroup := router.Group("")
	{
		userGroup.POST("/organization", middleware.Auth(),handlers.CreateOrganizationHandler)
		userGroup.POST("/organization/:organization_id/invite", middleware.Auth(), handlers.InviteUserHandler)
		userGroup.GET("/organization/:organization_id", middleware.Auth(), handlers.GetOrganizationByIDHandler)
		userGroup.PUT("/organization/:organization_id", middleware.Auth(), handlers.UpdateOrganizationHandler)
		userGroup.DELETE("/organization/:organization_id", middleware.Auth(), handlers.DeleteOrganizationHandler)
	}
}
