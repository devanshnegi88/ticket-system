package models

import "time"

// TicketStatus is a constrained string type representing the lifecycle
// state of a ticket.
type TicketStatus string

const (
	StatusOpen       TicketStatus = "open"
	StatusInProgress TicketStatus = "in_progress"
	StatusClosed     TicketStatus = "closed"
)

// IsValid reports whether the status is one of the three supported values.
func (s TicketStatus) IsValid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusClosed:
		return true
	default:
		return false
	}
}

// Ticket represents a support ticket owned by exactly one user.
type Ticket struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title" gorm:"not null"`
	Description string       `json:"description"`
	Status      TicketStatus `json:"status" gorm:"not null;default:open"`
	UserID      uint         `json:"user_id" gorm:"not null;index"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
