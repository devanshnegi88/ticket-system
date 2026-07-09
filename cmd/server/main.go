package main

import (
	"log"

	"ticket-system/configs"
	"ticket-system/internal/auth"
	"ticket-system/internal/database"
	"ticket-system/internal/handlers"
	"ticket-system/internal/repositories"
	"ticket-system/internal/services"
	"ticket-system/routes"
)

func main() {
	cfg := configs.Load()

	db := database.Connect(cfg.DatabasePath)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)

	userRepo := repositories.NewUserRepository(db)
	ticketRepo := repositories.NewTicketRepository(db)

	userService := services.NewUserService(userRepo, jwtManager)
	ticketService := services.NewTicketService(ticketRepo)

	authHandler := handlers.NewAuthHandler(userService)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	router := routes.Setup(authHandler, ticketHandler, jwtManager)

	log.Printf("starting ticket-system on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
