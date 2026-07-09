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
<<<<<<< HEAD
=======
//	GET    /              (HTML/CSS/JS UI)
//	GET    /css/*, /js/*  (static assets for the UI)
>>>>>>> 6d83ca9 (ticket system)
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
	router.Static("/css", "/web/css")
router.Static("/js", "/web/js")
router.StaticFile("/", "/web/index.html")

router.GET("/health", handlers.HealthCheck)




	// Serve the built-in HTML/CSS/JS UI. Static assets live under
	// ./web/css and ./web/js, and the SPA entry point is ./web/index.html,
	// served at "/" so the UI and API share an origin (no CORS needed).
	// router.Static("/css", "./web/css")
	// router.Static("/js", "./web/js")
	// router.StaticFile("/", "./web/index.html")

	// router.GET("/health", handlers.HealthCheck)


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
