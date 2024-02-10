package routes

import (
	"MoZaki-Organization-Manager/pkg/api/handlers"
	"MoZaki-Organization-Manager/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate)
	incomingRoutes.GET("/users", handlers.GetUsers)

}
