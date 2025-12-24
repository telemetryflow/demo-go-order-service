// Package repository defines repository interfaces.
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/telemetryflow/order-service/internal/domain/entity"
)

// OrderitemRepository defines the repository interface for Orderitem
type OrderitemRepository interface {
	// Create creates a new orderitem
	Create(ctx context.Context, e *entity.Orderitem) error

	// FindByID finds a orderitem by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Orderitem, error)

	// FindAll finds all orderitems with pagination
	FindAll(ctx context.Context, offset, limit int) ([]entity.Orderitem, int64, error)

	// Update updates an existing orderitem
	Update(ctx context.Context, e *entity.Orderitem) error

	// Delete soft-deletes a orderitem by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// HardDelete permanently deletes a orderitem
	HardDelete(ctx context.Context, id uuid.UUID) error
}
