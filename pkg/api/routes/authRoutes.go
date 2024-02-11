package routes

import (
	"MoZaki-Organization-Manager/pkg/api/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	auth := incomingRoutes.Group("/")
	{
		auth.POST("signup", handlers.Signup())
		auth.POST("signin", handlers.Login())
		auth.POST("refresh-token", handlers.RefreshToken())
	}

}
