package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ticket-system/internal/auth"
)

// ContextUserIDKey is the gin.Context key under which the authenticated
// user's ID is stored after AuthRequired succeeds.
const ContextUserIDKey = "userID"

// AuthRequired parses and validates the "Authorization: Bearer <token>"
// header. On success it stores the authenticated user's ID in the
// request context. On failure it aborts the request with 401.
func AuthRequired(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be in the format 'Bearer <token>'"})
			return
		}

		claims, err := jwtManager.Verify(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}
