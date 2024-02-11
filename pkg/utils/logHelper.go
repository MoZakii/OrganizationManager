package utils

import (
	"github.com/gin-gonic/gin"
)

func HasError(err error, msg string, c *gin.Context, status int) bool {
	if err != nil {
		c.JSON(status, gin.H{"error": msg})
		return true
	}
	return false
}
