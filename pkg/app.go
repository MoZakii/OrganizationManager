package pkg

import (
	"MoZaki-Organization-Manager/pkg/api/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Function that handles the initialization of the project
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

	router.Run(":" + port)
}
