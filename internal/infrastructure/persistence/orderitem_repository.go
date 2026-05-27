// Package persistence implements OrderItem repository with GORM.
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

// orderitemRepository implements repository.OrderitemRepository using GORM
type orderitemRepository struct {
	db *gorm.DB
}

// NewOrderitemRepository creates a new Orderitem repository
func NewOrderitemRepository(db *gorm.DB) repository.OrderitemRepository {
	return &orderitemRepository{
		db: db,
	}
}

// Create creates a new orderitem
func (r *orderitemRepository) Create(ctx context.Context, orderitem *entity.Orderitem) error {
	return r.db.WithContext(ctx).Create(orderitem).Error
}

// FindByID retrieves an orderitem by ID
func (r *orderitemRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Orderitem, error) {
	var orderitem entity.Orderitem
	err := r.db.WithContext(ctx).First(&orderitem, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("orderitem not found")
		}
		return nil, err
	}
	return &orderitem, nil
}

// FindAll retrieves all orderitems with pagination
func (r *orderitemRepository) FindAll(ctx context.Context, offset, limit int) ([]entity.Orderitem, int64, error) {
	var orderitems []entity.Orderitem
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&entity.Orderitem{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orderitems).Error; err != nil {
		return nil, 0, err
	}

	return orderitems, total, nil
}

// Update updates an orderitem
func (r *orderitemRepository) Update(ctx context.Context, orderitem *entity.Orderitem) error {
	return r.db.WithContext(ctx).Save(orderitem).Error
}

// Delete soft-deletes an orderitem by ID
func (r *orderitemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Orderitem{}, "id = ?", id).Error
}

// HardDelete permanently deletes an orderitem by ID
func (r *orderitemRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.Orderitem{}, "id = ?", id).Error
}

// FindByOrderID finds all items for an order
func (r *orderitemRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]entity.Orderitem, error) {
	var items []entity.Orderitem
	err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// FindByProductID finds all items for a product
func (r *orderitemRepository) FindByProductID(ctx context.Context, productID uuid.UUID) ([]entity.Orderitem, error) {
	var items []entity.Orderitem
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// CreateBatch creates multiple orderitems in a single transaction
func (r *orderitemRepository) CreateBatch(ctx context.Context, items []entity.Orderitem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}

// DeleteByOrderID deletes all items for an order
func (r *orderitemRepository) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Orderitem{}, "order_id = ?", orderID).Error
}
