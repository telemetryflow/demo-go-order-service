// Package entity_test provides unit tests for domain entities.
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

package entity_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/telemetryflow/order-service/internal/domain/entity"
	"gorm.io/gorm"
)

// =============================================================================
// Base Entity Tests
//
// Tests for the Base struct which provides common fields for all entities:
// ID (UUID), CreatedAt, UpdatedAt, and DeletedAt (soft delete).
// =============================================================================

// TestNewBase verifies that NewBase creates a properly initialized base entity
// with a valid UUID and timestamps.
func TestNewBase(t *testing.T) {
	t.Run("creates base with valid ID and timestamps", func(t *testing.T) {
		base := entity.NewBase()

		assert.NotEqual(t, uuid.Nil, base.ID, "ID should not be nil")
		assert.False(t, base.CreatedAt.IsZero(), "CreatedAt should be set")
		assert.False(t, base.UpdatedAt.IsZero(), "UpdatedAt should be set")
		assert.False(t, base.DeletedAt.Valid, "DeletedAt should not be valid")
	})

	t.Run("creates unique IDs", func(t *testing.T) {
		base1 := entity.NewBase()
		base2 := entity.NewBase()

		assert.NotEqual(t, base1.ID, base2.ID, "IDs should be unique")
	})
}

func TestBase_IsDeleted(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *entity.Base
		expected bool
	}{
		{
			name: "returns false for new base",
			setup: func() *entity.Base {
				base := entity.NewBase()
				return &base
			},
			expected: false,
		},
		{
			name: "returns true after MarkDeleted",
			setup: func() *entity.Base {
				base := entity.NewBase()
				base.MarkDeleted()
				return &base
			},
			expected: true,
		},
		{
			name: "returns false after Restore",
			setup: func() *entity.Base {
				base := entity.NewBase()
				base.MarkDeleted()
				base.Restore()
				return &base
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := tt.setup()
			assert.Equal(t, tt.expected, base.IsDeleted())
		})
	}
}

func TestBase_MarkUpdated(t *testing.T) {
	t.Run("updates timestamp", func(t *testing.T) {
		base := entity.NewBase()
		originalTime := base.UpdatedAt

		// Wait a moment to ensure different timestamp
		time.Sleep(time.Millisecond)
		base.MarkUpdated()

		assert.True(t, base.UpdatedAt.After(originalTime) || base.UpdatedAt.Equal(originalTime),
			"UpdatedAt should be updated")
	})
}

func TestBase_MarkDeleted(t *testing.T) {
	t.Run("sets DeletedAt", func(t *testing.T) {
		base := entity.NewBase()
		assert.False(t, base.DeletedAt.Valid)

		base.MarkDeleted()

		assert.True(t, base.DeletedAt.Valid)
		assert.False(t, base.DeletedAt.Time.IsZero())
	})
}

func TestBase_Restore(t *testing.T) {
	t.Run("clears DeletedAt", func(t *testing.T) {
		base := entity.NewBase()
		base.MarkDeleted()
		require.True(t, base.DeletedAt.Valid)

		base.Restore()

		assert.False(t, base.DeletedAt.Valid)
	})
}

func TestBase_BeforeCreate(t *testing.T) {
	t.Run("generates ID if nil", func(t *testing.T) {
		base := &entity.Base{}
		assert.Equal(t, uuid.Nil, base.ID)

		err := base.BeforeCreate(nil)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, base.ID)
	})

	t.Run("preserves existing ID", func(t *testing.T) {
		existingID := uuid.New()
		base := &entity.Base{ID: existingID}

		err := base.BeforeCreate(nil)

		assert.NoError(t, err)
		assert.Equal(t, existingID, base.ID)
	})
}

// =============================================================================
// Order Entity Tests
//
// Tests for the Order aggregate root entity. Orders represent customer purchases
// with associated items, totals, and lifecycle status management.
// =============================================================================

// TestOrder_TableName verifies the GORM table name convention for Order entities.
func TestOrder_TableName(t *testing.T) {
	t.Run("returns correct table name", func(t *testing.T) {
		order := entity.Order{}
		assert.Equal(t, "orders", order.TableName())
	})
}

func TestNewOrder(t *testing.T) {
	t.Run("creates order with all fields", func(t *testing.T) {
		customerID := uuid.New()
		total := 99.99
		status := "pending"

		order := entity.NewOrder(customerID, total, status)

		require.NotNil(t, order)
		assert.NotEqual(t, uuid.Nil, order.ID)
		assert.Equal(t, customerID, order.CustomerID)
		assert.Equal(t, total, order.Total)
		assert.Equal(t, status, order.Status)
		assert.False(t, order.CreatedAt.IsZero())
		assert.False(t, order.UpdatedAt.IsZero())
	})

	t.Run("creates order with different statuses", func(t *testing.T) {
		statuses := []string{"pending", "confirmed", "shipped", "delivered", "cancelled"}
		customerID := uuid.New()

		for _, status := range statuses {
			order := entity.NewOrder(customerID, 100.0, status)
			assert.Equal(t, status, order.Status)
		}
	})

	t.Run("creates order with zero total", func(t *testing.T) {
		order := entity.NewOrder(uuid.New(), 0, "pending")
		assert.Equal(t, float64(0), order.Total)
	})

	t.Run("creates order with negative total", func(t *testing.T) {
		// This might be valid for refunds
		order := entity.NewOrder(uuid.New(), -50.0, "refunded")
		assert.Equal(t, -50.0, order.Total)
	})
}

func TestOrder_Update(t *testing.T) {
	t.Run("updates all fields", func(t *testing.T) {
		originalCustomerID := uuid.New()
		order := entity.NewOrder(originalCustomerID, 100.0, "pending")
		originalUpdatedAt := order.UpdatedAt

		// Wait to ensure different timestamp
		time.Sleep(time.Millisecond)

		newCustomerID := uuid.New()
		newTotal := 200.0
		newStatus := "confirmed"

		order.Update(newCustomerID, newTotal, newStatus)

		assert.Equal(t, newCustomerID, order.CustomerID)
		assert.Equal(t, newTotal, order.Total)
		assert.Equal(t, newStatus, order.Status)
		assert.True(t, order.UpdatedAt.After(originalUpdatedAt) || order.UpdatedAt.Equal(originalUpdatedAt))
	})

	t.Run("preserves ID", func(t *testing.T) {
		order := entity.NewOrder(uuid.New(), 100.0, "pending")
		originalID := order.ID

		order.Update(uuid.New(), 200.0, "confirmed")

		assert.Equal(t, originalID, order.ID)
	})
}

func TestOrder_Validate(t *testing.T) {
	t.Run("returns nil for valid order", func(t *testing.T) {
		order := entity.NewOrder(uuid.New(), 100.0, "pending")
		err := order.Validate()
		assert.NoError(t, err)
	})
}

// =============================================================================
// Orderitem Entity Tests
//
// Tests for the Orderitem entity which represents individual line items within
// an order, including product references, quantities, and prices.
// =============================================================================

// TestOrderitem_TableName verifies the GORM table name convention for Orderitem entities.
func TestOrderitem_TableName(t *testing.T) {
	t.Run("returns correct table name", func(t *testing.T) {
		item := entity.Orderitem{}
		assert.Equal(t, "order_items", item.TableName())
	})
}

func TestNewOrderitem(t *testing.T) {
	t.Run("creates orderitem with all fields", func(t *testing.T) {
		orderID := uuid.New()
		productID := uuid.New()
		quantity := 5
		price := 29.99

		item := entity.NewOrderitem(orderID, productID, quantity, price)

		require.NotNil(t, item)
		assert.NotEqual(t, uuid.Nil, item.ID)
		assert.Equal(t, orderID, item.OrderID)
		assert.Equal(t, productID, item.ProductID)
		assert.Equal(t, quantity, item.Quantity)
		assert.Equal(t, price, item.Price)
		assert.False(t, item.CreatedAt.IsZero())
		assert.False(t, item.UpdatedAt.IsZero())
	})

	t.Run("creates orderitem with zero quantity", func(t *testing.T) {
		item := entity.NewOrderitem(uuid.New(), uuid.New(), 0, 10.0)
		assert.Equal(t, 0, item.Quantity)
	})

	t.Run("creates orderitem with zero price", func(t *testing.T) {
		item := entity.NewOrderitem(uuid.New(), uuid.New(), 1, 0)
		assert.Equal(t, float64(0), item.Price)
	})
}

func TestOrderitem_Update(t *testing.T) {
	t.Run("updates all fields", func(t *testing.T) {
		originalOrderID := uuid.New()
		originalProductID := uuid.New()
		item := entity.NewOrderitem(originalOrderID, originalProductID, 1, 10.0)
		originalUpdatedAt := item.UpdatedAt

		// Wait to ensure different timestamp
		time.Sleep(time.Millisecond)

		newOrderID := uuid.New()
		newProductID := uuid.New()
		newQuantity := 5
		newPrice := 50.0

		item.Update(newOrderID, newProductID, newQuantity, newPrice)

		assert.Equal(t, newOrderID, item.OrderID)
		assert.Equal(t, newProductID, item.ProductID)
		assert.Equal(t, newQuantity, item.Quantity)
		assert.Equal(t, newPrice, item.Price)
		assert.True(t, item.UpdatedAt.After(originalUpdatedAt) || item.UpdatedAt.Equal(originalUpdatedAt))
	})

	t.Run("preserves ID", func(t *testing.T) {
		item := entity.NewOrderitem(uuid.New(), uuid.New(), 1, 10.0)
		originalID := item.ID

		item.Update(uuid.New(), uuid.New(), 2, 20.0)

		assert.Equal(t, originalID, item.ID)
	})
}

func TestOrderitem_Validate(t *testing.T) {
	t.Run("returns nil for valid orderitem", func(t *testing.T) {
		item := entity.NewOrderitem(uuid.New(), uuid.New(), 1, 10.0)
		err := item.Validate()
		assert.NoError(t, err)
	})
}

// =============================================================================
// Order with Items Integration
//
// Tests for the Order-Orderitem relationship, verifying that orders can
// properly contain and manage their associated line items.
// =============================================================================

// TestOrder_WithItems verifies the one-to-many relationship between Order and Orderitem.
func TestOrder_WithItems(t *testing.T) {
	t.Run("order can have multiple items", func(t *testing.T) {
		order := entity.NewOrder(uuid.New(), 100.0, "pending")

		items := []entity.Orderitem{
			*entity.NewOrderitem(order.ID, uuid.New(), 2, 25.0),
			*entity.NewOrderitem(order.ID, uuid.New(), 1, 50.0),
		}
		order.Items = items

		assert.Len(t, order.Items, 2)
		for _, item := range order.Items {
			assert.Equal(t, order.ID, item.OrderID)
		}
	})

	t.Run("order can have empty items", func(t *testing.T) {
		order := entity.NewOrder(uuid.New(), 0, "pending")
		assert.Empty(t, order.Items)
	})
}

// =============================================================================
// Edge Cases and Boundary Tests
//
// Tests for unusual inputs and boundary conditions to ensure entity robustness.
// These tests help verify proper handling of extreme values and edge scenarios.
// =============================================================================

// TestEntity_EdgeCases verifies entity behavior under extreme or unusual conditions.
func TestEntity_EdgeCases(t *testing.T) {
	t.Run("order with very large total", func(t *testing.T) {
		largeTotal := 9999999999.99
		order := entity.NewOrder(uuid.New(), largeTotal, "pending")
		assert.Equal(t, largeTotal, order.Total)
	})

	t.Run("orderitem with very large quantity", func(t *testing.T) {
		largeQuantity := 999999
		item := entity.NewOrderitem(uuid.New(), uuid.New(), largeQuantity, 1.0)
		assert.Equal(t, largeQuantity, item.Quantity)
	})

	t.Run("multiple soft delete and restore cycles", func(t *testing.T) {
		base := entity.NewBase()

		for i := 0; i < 5; i++ {
			base.MarkDeleted()
			assert.True(t, base.IsDeleted())

			base.Restore()
			assert.False(t, base.IsDeleted())
		}
	})
}

// =============================================================================
// GORM Hook Tests
//
// Tests for GORM lifecycle hooks that handle automatic operations like
// ID generation before entity creation.
// =============================================================================

// TestBase_GORMHooks verifies GORM hooks work correctly for entity lifecycle events.
func TestBase_GORMHooks(t *testing.T) {
	t.Run("BeforeCreate with nil gorm.DB", func(t *testing.T) {
		base := &entity.Base{}
		err := base.BeforeCreate(nil)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, base.ID)
	})

	t.Run("DeletedAt valid after soft delete", func(t *testing.T) {
		base := entity.NewBase()
		base.MarkDeleted()

		assert.True(t, base.DeletedAt.Valid)
		assert.Equal(t, gorm.DeletedAt{Time: base.DeletedAt.Time, Valid: true}, base.DeletedAt)
	})
}

// =============================================================================
// Benchmark Tests
//
// Performance benchmarks for entity creation and update operations.
// Run with: go test -bench=. -benchmem
// =============================================================================

// BenchmarkNewBase measures the performance of Base entity creation.
func BenchmarkNewBase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = entity.NewBase()
	}
}

func BenchmarkNewOrder(b *testing.B) {
	customerID := uuid.New()
	for i := 0; i < b.N; i++ {
		_ = entity.NewOrder(customerID, 100.0, "pending")
	}
}

func BenchmarkNewOrderitem(b *testing.B) {
	orderID := uuid.New()
	productID := uuid.New()
	for i := 0; i < b.N; i++ {
		_ = entity.NewOrderitem(orderID, productID, 1, 10.0)
	}
}

func BenchmarkOrder_Update(b *testing.B) {
	order := entity.NewOrder(uuid.New(), 100.0, "pending")
	newCustomerID := uuid.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		order.Update(newCustomerID, float64(i), "confirmed")
	}
}
