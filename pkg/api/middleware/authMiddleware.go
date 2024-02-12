package middleware

import (
	"MoZaki-Organization-Manager/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Function used as middleware authentication to check on authorization token and proceed with the request if valid
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(clientToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user_email", claims["Email"])
		c.Set("user_name", claims["Name"])
		c.Set("user_id", claims["Uid"])
	}

}
