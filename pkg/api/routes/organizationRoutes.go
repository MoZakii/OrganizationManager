package routes

import (
	"MoZaki-Organization-Manager/pkg/api/handlers"
	"MoZaki-Organization-Manager/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func OrganizationRoutes(incomingRoutes *gin.Engine) {
	org := incomingRoutes.Group("/organization", middleware.Authenticate())
	{
		org.POST("", handlers.CreateOrganization())
		org.GET("/:organization_id", handlers.ReadOrganization())
		org.GET("", handlers.ReadAllOrganizations())
		org.PUT("/:organization_id", handlers.UpdateOrganization())
		org.DELETE("/:organization_id", handlers.DeleteOrganization())
		org.POST("/:organization_id/invite", handlers.InviteUserToOrganization())
	}
}
