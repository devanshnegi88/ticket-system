package repositories

import (
	"errors"

	"gorm.io/gorm"

	"ticket-system/internal/models"
)

// TicketRepository is the data-access layer for Ticket records.
type TicketRepository struct {
	db *gorm.DB
}

// NewTicketRepository builds a TicketRepository backed by the given DB
// connection.
func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

// Create inserts a new ticket row.
func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

// FindAllByUser returns every ticket owned by userID, most recent first.
func (r *TicketRepository) FindAllByUser(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

// FindByID returns a ticket by its primary key regardless of owner; the
// caller is responsible for enforcing ownership.
func (r *TicketRepository) FindByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.First(&ticket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &ticket, nil
}

// UpdateStatus persists a new status value for an already-loaded ticket.
func (r *TicketRepository) UpdateStatus(ticket *models.Ticket, status models.TicketStatus) error {
	if err := r.db.Model(ticket).Update("status", status).Error; err != nil {
		return err
	}
	ticket.Status = status
	return nil
}
