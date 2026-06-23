// Package entity contains the Order domain entity.
//
// TelemetryFlow Order Service - AI-Powered Observability & Incident Response Management (IRM) Platform
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

// Order represents the order domain entity
type Order struct {
	Base
	CustomerID uuid.UUID   `json:"customer_id" gorm:"type:uuid;not null;index"`
	Total      float64     `json:"total" gorm:"type:decimal(15,2);not null;default:0"`
	Status     string      `json:"status" gorm:"type:varchar(50);not null;default:'pending';index"`
	Items      []Orderitem `json:"items,omitempty" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName returns the table name for GORM
func (Order) TableName() string {
	return "orders"
}

// NewOrder creates a new Order entity
func NewOrder(customerID uuid.UUID, total float64, status string) *Order {
	return &Order{
		Base:       NewBase(),
		CustomerID: customerID,
		Total:      total,
		Status:     status,
	}
}

// Update updates the order fields
func (e *Order) Update(customerID uuid.UUID, total float64, status string) {
	e.CustomerID = customerID
	e.Total = total
	e.Status = status
	e.MarkUpdated()
}

// Validate validates the entity
func (e *Order) Validate() error {
	// Add validation logic here
	return nil
}
