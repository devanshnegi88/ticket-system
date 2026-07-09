package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/services"
	"ticket-system/internal/utils"
)

// TicketHandler exposes the /tickets endpoint group.
type TicketHandler struct {
	ticketService *services.TicketService
}

// NewTicketHandler builds a TicketHandler.
func NewTicketHandler(ticketService *services.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

type createTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type updateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// currentUserID extracts the authenticated user's ID set by
// middleware.AuthRequired. It should never be missing on a protected
// route, but is checked defensively.
func currentUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := val.(uint)
	return id, ok
}

// Create handles POST /tickets.
func (h *TicketHandler) Create(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	ticket, err := h.ticketService.Create(userID, req.Title, req.Description)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to create ticket")
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// List handles GET /tickets.
func (h *TicketHandler) List(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	tickets, err := h.ticketService.ListForUser(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "failed to list tickets")
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// GetByID handles GET /tickets/{id}.
func (h *TicketHandler) GetByID(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ticket id")
		return
	}

	ticket, err := h.ticketService.GetOwned(userID, uint(id))
	if err != nil {
		writeTicketServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// UpdateStatus handles PATCH /tickets/{id}/status.
func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid ticket id")
		return
	}

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	ticket, err := h.ticketService.UpdateStatus(userID, uint(id), models.TicketStatus(req.Status))
	if err != nil {
		writeTicketServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// writeTicketServiceError maps known service-layer errors to the
// appropriate HTTP status code and JSON error body.
func writeTicketServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrTicketNotFound):
		utils.ErrorResponse(c, http.StatusNotFound, "ticket not found")
	case errors.Is(err, services.ErrForbidden):
		utils.ErrorResponse(c, http.StatusForbidden, "you do not have access to this ticket")
	case errors.Is(err, services.ErrInvalidStatus):
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid status value")
	case errors.Is(err, services.ErrInvalidTransition):
		utils.ErrorResponse(c, http.StatusConflict, "invalid status transition")
	default:
		utils.ErrorResponse(c, http.StatusInternalServerError, "unexpected error")
	}
}
