// Package repository defines repository interfaces.
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/telemetryflow/order-service/internal/domain/entity"
)

// OrderRepository defines the repository interface for Order
type OrderRepository interface {
	// Create creates a new order
	Create(ctx context.Context, e *entity.Order) error

	// FindByID finds a order by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)

	// FindAll finds all orders with pagination
	FindAll(ctx context.Context, offset, limit int) ([]entity.Order, int64, error)

	// Update updates an existing order
	Update(ctx context.Context, e *entity.Order) error

	// Delete soft-deletes a order by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// HardDelete permanently deletes a order
	HardDelete(ctx context.Context, id uuid.UUID) error

	// FindByStatus finds orders by Status
	FindByStatus(ctx context.Context, status string) ([]entity.Order, error)
}
