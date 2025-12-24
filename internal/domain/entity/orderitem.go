// Package entity contains domain entities.
package entity

import (
	"github.com/google/uuid"
)

// Orderitem represents the orderitem domain entity
type Orderitem struct {
	Base
	OrderId   uuid.UUID `json:"orderId" db:"order_id"`
	ProductId uuid.UUID `json:"productId" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"`
}

// NewOrderitem creates a new Orderitem entity
func NewOrderitem(orderId uuid.UUID, productId uuid.UUID, quantity int, price float64) *Orderitem {
	return &Orderitem{
		Base:      NewBase(),
		OrderId:   orderId,
		ProductId: productId,
		Quantity:  quantity,
		Price:     price,
	}
}

// Update updates the orderitem fields
func (e *Orderitem) Update(orderId uuid.UUID, productId uuid.UUID, quantity int, price float64) {
	e.OrderId = orderId
	e.ProductId = productId
	e.Quantity = quantity
	e.Price = price
	e.MarkUpdated()
}

// Validate validates the entity
func (e *Orderitem) Validate() error {
	// Add validation logic here
	return nil
}
