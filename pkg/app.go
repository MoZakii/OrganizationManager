package pkg

import (
	"MoZaki-Organization-Manager/pkg/api/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.Default()

	routes.OrganizationRoutes(router)
	routes.AuthRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
