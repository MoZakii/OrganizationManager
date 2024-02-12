package routes

import (
	"MoZaki-Organization-Manager/pkg/api/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", handlers.Signup())
	incomingRoutes.POST("/signin", handlers.Login())
	incomingRoutes.POST("/refresh-token", handlers.RefreshToken())
}
