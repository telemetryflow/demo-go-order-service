// Package handler_test provides unit tests for application handlers.
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

package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/telemetryflow/order-service/internal/application/command"
	"github.com/telemetryflow/order-service/internal/application/handler"
	"github.com/telemetryflow/order-service/internal/application/query"
	"github.com/telemetryflow/order-service/internal/domain/entity"
)

// =============================================================================
// Mock Order Repository
//
// MockOrderRepository provides a test double for the OrderRepository interface,
// allowing tests to control repository behavior and verify interactions.
// =============================================================================

// MockOrderRepository implements repository.OrderRepository for testing.
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, e *entity.Order) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *MockOrderRepository) FindAll(ctx context.Context, offset, limit int) ([]entity.Order, int64, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]entity.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) Update(ctx context.Context, e *entity.Order) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByStatus(ctx context.Context, status string) ([]entity.Order, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockOrderRepository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]entity.Order, error) {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockOrderRepository) FindWithItems(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

// =============================================================================
// Order Command Handler Tests
//
// Tests for OrderCommandHandler which processes write operations (Create,
// Update, Delete) following the Command side of CQRS pattern.
// =============================================================================

// TestNewOrderCommandHandler verifies handler construction.
func TestNewOrderCommandHandler(t *testing.T) {
	t.Run("creates handler with repository", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderCommandHandler(repo)

		require.NotNil(t, h)
	})
}

func TestOrderCommandHandler_HandleOrderCreate(t *testing.T) {
	t.Run("successfully creates order", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderCommandHandler(repo)

		cmd := &command.CreateOrderCommand{
			CustomerID: uuid.New(),
			Total:      100.0,
			Status:     "pending",
		}

		repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

		err := h.HandleOrderCreate(context.Background(), cmd)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderCommandHandler(repo)

		cmd := &command.CreateOrderCommand{
			CustomerID: uuid.New(),
			Total:      100.0,
			Status:     "pending",
		}

		expectedErr := errors.New("database error")
		repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(expectedErr)

		err := h.HandleOrderCreate(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestOrderCommandHandler_HandleOrderUpdate(t *testing.T) {
	t.Run("successfully updates order", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderCommandHandler(repo)

		cmd := &command.UpdateOrderCommand{
			ID:         uuid.New(),
			CustomerID: uuid.New(),
			Total:      200.0,
			Status:     "confirmed",
		}

		repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

		err := h.HandleOrderUpdate(context.Background(), cmd)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderCommandHandler(repo)

		cmd := &command.UpdateOrderCommand{
			ID:         uuid.New(),
			CustomerID: uuid.New(),
			Total:      200.0,
			Status:     "confirmed",
		}

		expectedErr := errors.New("update failed")
		repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(expectedErr)

		err := h.HandleOrderUpdate(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestOrderCommandHandler_HandleOrderDelete(t *testing.T) {
	t.Run("successfully deletes order", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderCommandHandler(repo)

		orderID := uuid.New()
		cmd := &command.DeleteOrderCommand{ID: orderID}

		repo.On("Delete", mock.Anything, orderID).Return(nil)

		err := h.HandleOrderDelete(context.Background(), cmd)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderCommandHandler(repo)

		orderID := uuid.New()
		cmd := &command.DeleteOrderCommand{ID: orderID}

		expectedErr := errors.New("delete failed")
		repo.On("Delete", mock.Anything, orderID).Return(expectedErr)

		err := h.HandleOrderDelete(context.Background(), cmd)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

// =============================================================================
// Order Query Handler Tests
//
// Tests for OrderQueryHandler which processes read operations (GetByID,
// GetAll) following the Query side of CQRS pattern.
// =============================================================================

// TestNewOrderQueryHandler verifies handler construction.
func TestNewOrderQueryHandler(t *testing.T) {
	t.Run("creates handler with repository", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderQueryHandler(repo)

		require.NotNil(t, h)
	})
}

func TestOrderQueryHandler_HandleOrderGetByID(t *testing.T) {
	t.Run("successfully gets order by ID", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderQueryHandler(repo)

		orderID := uuid.New()
		customerID := uuid.New()
		expectedOrder := entity.NewOrder(customerID, 100.0, "pending")
		expectedOrder.ID = orderID

		qry := &query.GetOrderByIDQuery{ID: orderID}

		repo.On("FindByID", mock.Anything, orderID).Return(expectedOrder, nil)

		result, err := h.HandleOrderGetByID(context.Background(), qry)

		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, orderID, result.ID)
		assert.Equal(t, customerID, result.CustomerID)
		assert.Equal(t, 100.0, result.Total)
		assert.Equal(t, "pending", result.Status)
		repo.AssertExpectations(t)
	})

	t.Run("returns error when order not found", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderQueryHandler(repo)

		orderID := uuid.New()
		qry := &query.GetOrderByIDQuery{ID: orderID}

		expectedErr := errors.New("order not found")
		repo.On("FindByID", mock.Anything, orderID).Return(nil, expectedErr)

		result, err := h.HandleOrderGetByID(context.Background(), qry)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}

func TestOrderQueryHandler_HandleOrderGetAll(t *testing.T) {
	t.Run("successfully gets all orders", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderQueryHandler(repo)

		orders := []entity.Order{
			*entity.NewOrder(uuid.New(), 100.0, "pending"),
			*entity.NewOrder(uuid.New(), 200.0, "confirmed"),
		}

		qry := &query.GetAllOrdersQuery{Offset: 0, Limit: 10}

		repo.On("FindAll", mock.Anything, 0, 10).Return(orders, int64(2), nil)

		result, err := h.HandleOrderGetAll(context.Background(), qry)

		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 0, result.Offset)
		assert.Equal(t, 10, result.Limit)
		repo.AssertExpectations(t)
	})

	t.Run("returns empty list when no orders", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderQueryHandler(repo)

		qry := &query.GetAllOrdersQuery{Offset: 0, Limit: 10}

		repo.On("FindAll", mock.Anything, 0, 10).Return([]entity.Order{}, int64(0), nil)

		result, err := h.HandleOrderGetAll(context.Background(), qry)

		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Data, 0)
		assert.Equal(t, 0, result.Total)
		repo.AssertExpectations(t)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderQueryHandler(repo)

		qry := &query.GetAllOrdersQuery{Offset: 0, Limit: 10}

		expectedErr := errors.New("database error")
		repo.On("FindAll", mock.Anything, 0, 10).Return(nil, int64(0), expectedErr)

		result, err := h.HandleOrderGetAll(context.Background(), qry)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})

	t.Run("handles pagination correctly", func(t *testing.T) {
		repo := new(MockOrderRepository)
		h := handler.NewOrderQueryHandler(repo)

		orders := []entity.Order{
			*entity.NewOrder(uuid.New(), 300.0, "shipped"),
		}

		qry := &query.GetAllOrdersQuery{Offset: 10, Limit: 5}

		repo.On("FindAll", mock.Anything, 10, 5).Return(orders, int64(15), nil)

		result, err := h.HandleOrderGetAll(context.Background(), qry)

		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Data, 1)
		assert.Equal(t, 15, result.Total)
		assert.Equal(t, 10, result.Offset)
		assert.Equal(t, 5, result.Limit)
		repo.AssertExpectations(t)
	})
}

// =============================================================================
// Integration-style Handler Tests
//
// Tests that verify the full CRUD workflow through both command and query
// handlers, simulating real-world usage patterns.
// =============================================================================

// TestOrderHandler_FullWorkflow verifies the complete create, read, update, delete cycle.
func TestOrderHandler_FullWorkflow(t *testing.T) {
	t.Run("create, read, update, delete workflow", func(t *testing.T) {
		repo := new(MockOrderRepository)
		cmdHandler := handler.NewOrderCommandHandler(repo)
		qryHandler := handler.NewOrderQueryHandler(repo)

		orderID := uuid.New()
		customerID := uuid.New()

		// Create
		createCmd := &command.CreateOrderCommand{
			CustomerID: customerID,
			Total:      100.0,
			Status:     "pending",
		}
		repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil).Once()

		err := cmdHandler.HandleOrderCreate(context.Background(), createCmd)
		assert.NoError(t, err)

		// Read
		order := entity.NewOrder(customerID, 100.0, "pending")
		order.ID = orderID
		repo.On("FindByID", mock.Anything, orderID).Return(order, nil).Once()

		qry := &query.GetOrderByIDQuery{ID: orderID}
		result, err := qryHandler.HandleOrderGetByID(context.Background(), qry)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		// Update
		updateCmd := &command.UpdateOrderCommand{
			ID:         orderID,
			CustomerID: customerID,
			Total:      150.0,
			Status:     "confirmed",
		}
		repo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil).Once()

		err = cmdHandler.HandleOrderUpdate(context.Background(), updateCmd)
		assert.NoError(t, err)

		// Delete
		deleteCmd := &command.DeleteOrderCommand{ID: orderID}
		repo.On("Delete", mock.Anything, orderID).Return(nil).Once()

		err = cmdHandler.HandleOrderDelete(context.Background(), deleteCmd)
		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})
}

// =============================================================================
// Benchmark Tests
//
// Performance benchmarks for handler operations.
// Run with: go test -bench=. -benchmem
// =============================================================================

// BenchmarkOrderCommandHandler_HandleOrderCreate measures create handler performance.
func BenchmarkOrderCommandHandler_HandleOrderCreate(b *testing.B) {
	repo := new(MockOrderRepository)
	h := handler.NewOrderCommandHandler(repo)

	cmd := &command.CreateOrderCommand{
		CustomerID: uuid.New(),
		Total:      100.0,
		Status:     "pending",
	}

	repo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = h.HandleOrderCreate(context.Background(), cmd)
	}
}

func BenchmarkOrderQueryHandler_HandleOrderGetByID(b *testing.B) {
	repo := new(MockOrderRepository)
	h := handler.NewOrderQueryHandler(repo)

	orderID := uuid.New()
	order := entity.NewOrder(uuid.New(), 100.0, "pending")
	order.ID = orderID

	qry := &query.GetOrderByIDQuery{ID: orderID}

	repo.On("FindByID", mock.Anything, orderID).Return(order, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = h.HandleOrderGetByID(context.Background(), qry)
	}
}
