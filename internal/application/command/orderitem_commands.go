// Package command contains CQRS commands for Orderitem.
package command

import (
	"github.com/google/uuid"
	"github.com/telemetryflow/order-service/internal/domain/entity"
)

// CreateOrderitemCommand represents the create orderitem command
type CreateOrderitemCommand struct {
	OrderId   uuid.UUID `json:"orderId" validate:"required"`
	ProductId uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
	Price     float64   `json:"price" validate:"required"`
}

// Validate validates the create command
func (c *CreateOrderitemCommand) Validate() error {
	// Add validation logic
	return nil
}

// ToEntity converts the command to an entity
func (c *CreateOrderitemCommand) ToEntity() *entity.Orderitem {
	return entity.NewOrderitem(c.OrderId, c.ProductId, c.Quantity, c.Price)
}

// UpdateOrderitemCommand represents the update orderitem command
type UpdateOrderitemCommand struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	OrderId   uuid.UUID `json:"orderId" validate:"required"`
	ProductId uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
	Price     float64   `json:"price" validate:"required"`
}

// Validate validates the update command
func (c *UpdateOrderitemCommand) Validate() error {
	if c.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}

// ToEntity converts the command to an entity
func (c *UpdateOrderitemCommand) ToEntity() *entity.Orderitem {
	e := entity.NewOrderitem(c.OrderId, c.ProductId, c.Quantity, c.Price)
	e.ID = c.ID
	return e
}

// DeleteOrderitemCommand represents the delete orderitem command
type DeleteOrderitemCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// Validate validates the delete command
func (c *DeleteOrderitemCommand) Validate() error {
	if c.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}
