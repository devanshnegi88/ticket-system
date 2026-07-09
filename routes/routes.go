package routes

import (
	"github.com/gin-gonic/gin"

	"ticket-system/internal/auth"
	"ticket-system/internal/handlers"
	"ticket-system/internal/middleware"
)

// Setup builds the gin.Engine with all routes registered, matching the
// endpoint contract required by the assignment exactly:
//
//	GET    /health
//	POST   /auth/register
//	POST   /auth/login
//	POST   /tickets
//	GET    /tickets
//	GET    /tickets/{id}
//	PATCH  /tickets/{id}/status
func Setup(
	authHandler *handlers.AuthHandler,
	ticketHandler *handlers.TicketHandler,
	jwtManager *auth.JWTManager,
) *gin.Engine {
	router := gin.Default()

	router.GET("/health", handlers.HealthCheck)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	ticketGroup := router.Group("/tickets")
	ticketGroup.Use(middleware.AuthRequired(jwtManager))
	{
		ticketGroup.POST("", ticketHandler.Create)
		ticketGroup.GET("", ticketHandler.List)
		ticketGroup.GET("/:id", ticketHandler.GetByID)
		ticketGroup.PATCH("/:id/status", ticketHandler.UpdateStatus)
	}

	return router
}
