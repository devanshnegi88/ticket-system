package services

import (
	"errors"

	"ticket-system/internal/models"
	"ticket-system/internal/repositories"
)

var (
	// ErrTicketNotFound is returned when no ticket exists with the given ID.
	ErrTicketNotFound = errors.New("ticket not found")
	// ErrForbidden is returned when a ticket exists but belongs to a
	// different user.
	ErrForbidden = errors.New("you do not have access to this ticket")
	// ErrInvalidStatus is returned when the requested status is not one
	// of open, in_progress, closed.
	ErrInvalidStatus = errors.New("invalid status value")
	// ErrInvalidTransition is returned when the requested status change
	// does not follow the required open -> in_progress -> closed flow,
	// including any attempt to reopen a closed ticket.
	ErrInvalidTransition = errors.New("invalid status transition")
)

// allowedTransitions maps a ticket's current status to the single
// next status it is permitted to move to. Statuses with no entry
// (i.e. StatusClosed) cannot transition anywhere.
var allowedTransitions = map[models.TicketStatus]models.TicketStatus{
	models.StatusOpen:       models.StatusInProgress,
	models.StatusInProgress: models.StatusClosed,
}

// TicketService contains the business logic for tickets, including the
// ownership enforcement and status-transition rules required by the
// assignment.
type TicketService struct {
	repo *repositories.TicketRepository
}

// NewTicketService builds a TicketService.
func NewTicketService(repo *repositories.TicketRepository) *TicketService {
	return &TicketService{repo: repo}
}

// Create makes a new ticket for userID, always starting in the "open"
// status.
func (s *TicketService) Create(userID uint, title, description string) (*models.Ticket, error) {
	ticket := &models.Ticket{
		Title:       title,
		Description: description,
		Status:      models.StatusOpen,
		UserID:      userID,
	}
	if err := s.repo.Create(ticket); err != nil {
		return nil, err
	}
	return ticket, nil
}

// ListForUser returns every ticket owned by userID.
func (s *TicketService) ListForUser(userID uint) ([]models.Ticket, error) {
	return s.repo.FindAllByUser(userID)
}

// GetOwned fetches a ticket by ID and verifies that it belongs to
// userID, returning ErrTicketNotFound or ErrForbidden as appropriate.
func (s *TicketService) GetOwned(userID, ticketID uint) (*models.Ticket, error) {
	ticket, err := s.repo.FindByID(ticketID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}
	if ticket.UserID != userID {
		return nil, ErrForbidden
	}
	return ticket, nil
}

// UpdateStatus validates and applies a status change to a ticket owned
// by userID. It enforces:
//   - the new status must be one of open/in_progress/closed
//   - ownership (via GetOwned)
//   - the strict sequential flow open -> in_progress -> closed
//   - a closed ticket can never move to any other status
func (s *TicketService) UpdateStatus(userID, ticketID uint, newStatus models.TicketStatus) (*models.Ticket, error) {
	if !newStatus.IsValid() {
		return nil, ErrInvalidStatus
	}

	ticket, err := s.GetOwned(userID, ticketID)
	if err != nil {
		return nil, err
	}

	next, allowed := allowedTransitions[ticket.Status]
	if !allowed || next != newStatus {
		return nil, ErrInvalidTransition
	}

	if err := s.repo.UpdateStatus(ticket, newStatus); err != nil {
		return nil, err
	}
	return ticket, nil
}
