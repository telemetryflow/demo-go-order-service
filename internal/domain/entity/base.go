// Package entity contains domain entities.
package entity

import (
	"time"

	"github.com/google/uuid"
)

// Base contains common fields for all entities
type Base struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// NewBase creates a new Base with generated ID and timestamps
func NewBase() Base {
	now := time.Now()
	return Base{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsDeleted returns true if the entity has been soft-deleted
func (b *Base) IsDeleted() bool {
	return b.DeletedAt != nil
}

// MarkUpdated updates the UpdatedAt timestamp
func (b *Base) MarkUpdated() {
	b.UpdatedAt = time.Now()
}

// MarkDeleted sets the DeletedAt timestamp for soft delete
func (b *Base) MarkDeleted() {
	now := time.Now()
	b.DeletedAt = &now
}

// Restore clears the DeletedAt timestamp
func (b *Base) Restore() {
	b.DeletedAt = nil
}
