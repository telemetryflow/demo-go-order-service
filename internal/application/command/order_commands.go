// Package command contains CQRS commands for Order.
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

package command

import (
	"github.com/google/uuid"
	"github.com/telemetryflow/order-service/internal/domain/entity"
)

// CreateOrderCommand represents the create order command
type CreateOrderCommand struct {
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Total      float64   `json:"total" validate:"required"`
	Status     string    `json:"status" validate:"required"`
}

// Validate validates the create command
func (c *CreateOrderCommand) Validate() error {
	// Add validation logic
	return nil
}

// ToEntity converts the command to an entity
func (c *CreateOrderCommand) ToEntity() *entity.Order {
	return entity.NewOrder(c.CustomerID, c.Total, c.Status)
}

// UpdateOrderCommand represents the update order command
type UpdateOrderCommand struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Total      float64   `json:"total" validate:"required"`
	Status     string    `json:"status" validate:"required"`
}

// Validate validates the update command
func (c *UpdateOrderCommand) Validate() error {
	if c.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}

// ToEntity converts the command to an entity
func (c *UpdateOrderCommand) ToEntity() *entity.Order {
	e := entity.NewOrder(c.CustomerID, c.Total, c.Status)
	e.ID = c.ID
	return e
}

// DeleteOrderCommand represents the delete order command
type DeleteOrderCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// Validate validates the delete command
func (c *DeleteOrderCommand) Validate() error {
	if c.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}
