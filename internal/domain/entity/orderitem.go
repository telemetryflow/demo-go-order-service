// Package entity contains the OrderItem domain entity.
//
// TelemetryFlow Order Service - Community Enterprise Observability Platform
// Copyright (c) 2024-2026 Telemetri Data Indonesia. All rights reserved.
// Open Source Software built by Telemetri Data Indonesia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entity

import (
	"github.com/google/uuid"
)

// Orderitem represents the orderitem domain entity
type Orderitem struct {
	Base
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Price     float64   `json:"price" gorm:"type:decimal(15,2);not null;default:0"`
}

// TableName returns the table name for GORM
func (Orderitem) TableName() string {
	return "order_items"
}

// NewOrderitem creates a new Orderitem entity
func NewOrderitem(orderID uuid.UUID, productID uuid.UUID, quantity int, price float64) *Orderitem {
	return &Orderitem{
		Base:      NewBase(),
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
	}
}

// Update updates the orderitem fields
func (e *Orderitem) Update(orderID uuid.UUID, productID uuid.UUID, quantity int, price float64) {
	e.OrderID = orderID
	e.ProductID = productID
	e.Quantity = quantity
	e.Price = price
	e.MarkUpdated()
}

// Validate validates the entity
func (e *Orderitem) Validate() error {
	// Add validation logic here
	return nil
}
