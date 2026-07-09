package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"ticket-system/internal/services"
	"ticket-system/internal/utils"
)

// AuthHandler exposes /auth/register and /auth/login.
type AuthHandler struct {
	userService *services.UserService
}

// NewAuthHandler builds an AuthHandler.
func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	user, err := h.userService.Register(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailExists) {
			utils.ErrorResponse(c, http.StatusConflict, "email already registered")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to register user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.ErrorResponse(c, http.StatusUnauthorized, "invalid email or password")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to login")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
