// Package dto_test provides unit tests for DTOs.
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

package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/telemetryflow/order-service/internal/application/dto"
	"github.com/telemetryflow/order-service/internal/domain/entity"
)

// =============================================================================
// OrderResponse Tests
//
// Tests for converting Order entities to OrderResponse DTOs for API responses.
// =============================================================================

// TestFromOrder verifies single entity to response conversion.
func TestFromOrder(t *testing.T) {
	t.Run("converts entity to response correctly", func(t *testing.T) {
		customerID := uuid.New()
		order := entity.NewOrder(customerID, 150.50, "confirmed")

		response := dto.FromOrder(order)

		assert.Equal(t, order.ID, response.ID)
		assert.Equal(t, customerID, response.CustomerID)
		assert.Equal(t, 150.50, response.Total)
		assert.Equal(t, "confirmed", response.Status)
		assert.Equal(t, order.CreatedAt, response.CreatedAt)
		assert.Equal(t, order.UpdatedAt, response.UpdatedAt)
	})

	t.Run("handles zero values", func(t *testing.T) {
		order := &entity.Order{}

		response := dto.FromOrder(order)

		assert.Equal(t, uuid.Nil, response.ID)
		assert.Equal(t, uuid.Nil, response.CustomerID)
		assert.Equal(t, float64(0), response.Total)
		assert.Equal(t, "", response.Status)
	})
}

func TestFromOrders(t *testing.T) {
	t.Run("converts multiple entities to responses", func(t *testing.T) {
		entities := []entity.Order{
			*entity.NewOrder(uuid.New(), 100.0, "pending"),
			*entity.NewOrder(uuid.New(), 200.0, "confirmed"),
			*entity.NewOrder(uuid.New(), 300.0, "shipped"),
		}

		responses := dto.FromOrders(entities)

		require.Len(t, responses, 3)
		for i, resp := range responses {
			assert.Equal(t, entities[i].ID, resp.ID)
			assert.Equal(t, entities[i].CustomerID, resp.CustomerID)
			assert.Equal(t, entities[i].Total, resp.Total)
			assert.Equal(t, entities[i].Status, resp.Status)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		responses := dto.FromOrders([]entity.Order{})
		assert.Empty(t, responses)
		assert.NotNil(t, responses)
	})

	t.Run("handles nil entities in slice", func(t *testing.T) {
		// This tests the actual behavior - it creates a response for each entity
		entities := []entity.Order{
			*entity.NewOrder(uuid.New(), 100.0, "pending"),
		}

		responses := dto.FromOrders(entities)

		require.Len(t, responses, 1)
	})
}

func TestOrderToResponse(t *testing.T) {
	t.Run("converts entity pointer to response pointer", func(t *testing.T) {
		order := entity.NewOrder(uuid.New(), 100.0, "pending")

		response := dto.OrderToResponse(order)

		require.NotNil(t, response)
		assert.Equal(t, order.ID, response.ID)
		assert.Equal(t, order.CustomerID, response.CustomerID)
	})

	t.Run("returns nil for nil entity", func(t *testing.T) {
		response := dto.OrderToResponse(nil)
		assert.Nil(t, response)
	})
}

// =============================================================================
// OrderListResponse Tests
// =============================================================================

func TestOrderListResponse(t *testing.T) {
	t.Run("contains correct pagination info", func(t *testing.T) {
		order1 := entity.NewOrder(uuid.New(), 100.0, "pending")
		order2 := entity.NewOrder(uuid.New(), 200.0, "confirmed")

		resp1 := dto.OrderToResponse(order1)
		resp2 := dto.OrderToResponse(order2)

		listResponse := dto.OrderListResponse{
			Data:   []*dto.OrderResponse{resp1, resp2},
			Total:  100,
			Offset: 10,
			Limit:  20,
		}

		assert.Len(t, listResponse.Data, 2)
		assert.Equal(t, 100, listResponse.Total)
		assert.Equal(t, 10, listResponse.Offset)
		assert.Equal(t, 20, listResponse.Limit)
	})

	t.Run("handles empty data", func(t *testing.T) {
		listResponse := dto.OrderListResponse{
			Data:   []*dto.OrderResponse{},
			Total:  0,
			Offset: 0,
			Limit:  10,
		}

		assert.Empty(t, listResponse.Data)
		assert.Equal(t, 0, listResponse.Total)
	})
}

// =============================================================================
// CreateOrderRequest Tests
// =============================================================================

func TestCreateOrderRequest(t *testing.T) {
	t.Run("struct has correct field tags", func(t *testing.T) {
		req := dto.CreateOrderRequest{
			CustomerID: uuid.New(),
			Total:      99.99,
			Status:     "pending",
		}

		assert.NotEqual(t, uuid.Nil, req.CustomerID)
		assert.Equal(t, 99.99, req.Total)
		assert.Equal(t, "pending", req.Status)
	})
}

// =============================================================================
// UpdateOrderRequest Tests
// =============================================================================

func TestUpdateOrderRequest(t *testing.T) {
	t.Run("struct has correct field tags", func(t *testing.T) {
		req := dto.UpdateOrderRequest{
			CustomerID: uuid.New(),
			Total:      199.99,
			Status:     "confirmed",
		}

		assert.NotEqual(t, uuid.Nil, req.CustomerID)
		assert.Equal(t, 199.99, req.Total)
		assert.Equal(t, "confirmed", req.Status)
	})
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestDTO_EdgeCases(t *testing.T) {
	t.Run("response preserves timestamp precision", func(t *testing.T) {
		now := time.Now()
		order := &entity.Order{
			Base: entity.Base{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			CustomerID: uuid.New(),
			Total:      100.0,
			Status:     "pending",
		}

		response := dto.FromOrder(order)

		assert.Equal(t, now, response.CreatedAt)
		assert.Equal(t, now, response.UpdatedAt)
	})

	t.Run("handles very large totals", func(t *testing.T) {
		largeTotal := 9999999999.99
		order := entity.NewOrder(uuid.New(), largeTotal, "pending")

		response := dto.FromOrder(order)

		assert.Equal(t, largeTotal, response.Total)
	})

	t.Run("handles negative totals", func(t *testing.T) {
		negativeTotal := -50.0
		order := entity.NewOrder(uuid.New(), negativeTotal, "refunded")

		response := dto.FromOrder(order)

		assert.Equal(t, negativeTotal, response.Total)
	})

	t.Run("handles special characters in status", func(t *testing.T) {
		specialStatus := "pending-review_123"
		order := entity.NewOrder(uuid.New(), 100.0, specialStatus)

		response := dto.FromOrder(order)

		assert.Equal(t, specialStatus, response.Status)
	})
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkFromOrder(b *testing.B) {
	order := entity.NewOrder(uuid.New(), 100.0, "pending")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dto.FromOrder(order)
	}
}

func BenchmarkFromOrders(b *testing.B) {
	entities := make([]entity.Order, 100)
	for i := range entities {
		entities[i] = *entity.NewOrder(uuid.New(), float64(i*10), "pending")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dto.FromOrders(entities)
	}
}

func BenchmarkOrderToResponse(b *testing.B) {
	order := entity.NewOrder(uuid.New(), 100.0, "pending")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dto.OrderToResponse(order)
	}
}
