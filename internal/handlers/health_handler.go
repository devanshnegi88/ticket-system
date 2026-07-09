package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck responds with the exact payload required by the
// assignment: {"status": "ok"}.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
