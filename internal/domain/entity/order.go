// Package entity contains domain entities.
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Order represents the order domain entity
type Order struct {
	Base
	CustomerId uuid.UUID `json:"customerId" db:"customer_id"`
	Total      float64   `json:"total" db:"total"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
}

// NewOrder creates a new Order entity
func NewOrder(customerId uuid.UUID, total float64, status string, createdAt time.Time) *Order {
	return &Order{
		Base:       NewBase(),
		CustomerId: customerId,
		Total:      total,
		Status:     status,
		CreatedAt:  createdAt,
	}
}

// Update updates the order fields
func (e *Order) Update(customerId uuid.UUID, total float64, status string, createdAt time.Time) {
	e.CustomerId = customerId
	e.Total = total
	e.Status = status
	e.CreatedAt = createdAt
	e.MarkUpdated()
}

// Validate validates the entity
func (e *Order) Validate() error {
	// Add validation logic here
	return nil
}
