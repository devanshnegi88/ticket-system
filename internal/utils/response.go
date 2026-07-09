package utils

import "github.com/gin-gonic/gin"

// ErrorResponse writes a standard {"error": "..."} JSON body with the
// given HTTP status code.
func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
