// Package repository defines the Order repository interface.
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

	// FindByCustomerID finds orders by customer ID
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]entity.Order, error)

	// FindWithItems finds an order with its items
	FindWithItems(ctx context.Context, id uuid.UUID) (*entity.Order, error)
}
