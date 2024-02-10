package routes

import (
	"MoZaki-Organization-Manager/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func OrganizationRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate)
	org := incomingRoutes.Group("/organization")
	{
		org.POST("/", handlers.CreateOrganization)
		org.GET("/:organization_id", handlers.ReadOrganization)
		org.GET("/", handlers.ReadAllOrganizations)
		org.PUT("/:organization_id", handlers.UpdateOrganization)
		org.DELETE("/:organization_id", handlers.DeleteOrganization)
		org.POST("/:organization_id/invite", handlers.InviteUserToOrganization)
	}
}
