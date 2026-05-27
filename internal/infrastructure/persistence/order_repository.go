// Package persistence implements Order repository with GORM.
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

package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/telemetryflow/order-service/internal/domain/entity"
	"github.com/telemetryflow/order-service/internal/domain/repository"
	"gorm.io/gorm"
)

// orderRepository implements repository.OrderRepository using GORM
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new Order repository
func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &orderRepository{
		db: db,
	}
}

// Create creates a new order
func (r *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// FindByID retrieves an order by ID
func (r *orderRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.WithContext(ctx).First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

// FindAll retrieves all orders with pagination
func (r *orderRepository) FindAll(ctx context.Context, offset, limit int) ([]entity.Order, int64, error) {
	var orders []entity.Order
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&entity.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// Update updates an order
func (r *orderRepository) Update(ctx context.Context, order *entity.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete soft-deletes an order by ID
func (r *orderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Order{}, "id = ?", id).Error
}

// HardDelete permanently deletes an order by ID
func (r *orderRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.Order{}, "id = ?", id).Error
}

// FindByStatus finds orders by status
func (r *orderRepository) FindByStatus(ctx context.Context, status string) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// FindByCustomerID finds orders by customer ID
func (r *orderRepository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.WithContext(ctx).
		Where("customer_id = ?", customerID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// FindWithItems retrieves an order with its items
func (r *orderRepository) FindWithItems(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}
