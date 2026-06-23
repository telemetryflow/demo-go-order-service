// Package http_test provides unit tests for HTTP handlers.
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

package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/telemetryflow/order-service/internal/application/command"
	"github.com/telemetryflow/order-service/internal/application/dto"
	apphandler "github.com/telemetryflow/order-service/internal/application/handler"
	"github.com/telemetryflow/order-service/internal/application/query"
	"github.com/telemetryflow/order-service/internal/domain/entity"
	httphandler "github.com/telemetryflow/order-service/internal/infrastructure/http/handler"
	"github.com/telemetryflow/order-service/pkg/validator"
)

// =============================================================================
// Mock Order Repository
//
// MockOrderRepository provides a test double for the OrderRepository interface.
// =============================================================================

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
// Mock Handlers for HTTP Handler Tests
// =============================================================================

type MockOrderCommandHandler struct {
	mock.Mock
}

func (m *MockOrderCommandHandler) HandleOrderCreate(ctx context.Context, cmd *command.CreateOrderCommand) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func (m *MockOrderCommandHandler) HandleOrderUpdate(ctx context.Context, cmd *command.UpdateOrderCommand) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

func (m *MockOrderCommandHandler) HandleOrderDelete(ctx context.Context, cmd *command.DeleteOrderCommand) error {
	args := m.Called(ctx, cmd)
	return args.Error(0)
}

type MockOrderQueryHandler struct {
	mock.Mock
}

func (m *MockOrderQueryHandler) HandleOrderGetByID(ctx context.Context, qry *query.GetOrderByIDQuery) (*dto.OrderResponse, error) {
	args := m.Called(ctx, qry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.OrderResponse), args.Error(1)
}

func (m *MockOrderQueryHandler) HandleOrderGetAll(ctx context.Context, qry *query.GetAllOrdersQuery) (*dto.OrderListResponse, error) {
	args := m.Called(ctx, qry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.OrderListResponse), args.Error(1)
}

// =============================================================================
// Order HTTP Handler Tests
// =============================================================================

func setupOrderHandlerTest() (*echo.Echo, *MockOrderRepository) {
	e := echo.New()
	e.Validator = validator.NewEchoValidator()
	mockRepo := new(MockOrderRepository)
	return e, mockRepo
}

func TestNewOrderHandler(t *testing.T) {
	t.Run("creates handler with dependencies", func(t *testing.T) {
		mockRepo := new(MockOrderRepository)
		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)

		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		require.NotNil(t, h)
	})
}

func TestOrderHandler_Create(t *testing.T) {
	t.Run("successfully creates order", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		customerID := uuid.New()
		reqBody := `{"customer_id":"` + customerID.String() + `","total":100.50,"status":"pending"}`

		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 400 for invalid JSON", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 400 for validation error", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		// Missing required fields
		reqBody := `{"total":100.50}`

		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 500 on repository error", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		customerID := uuid.New()
		reqBody := `{"customer_id":"` + customerID.String() + `","total":100.50,"status":"pending"}`

		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(errors.New("database error"))

		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestOrderHandler_GetByID(t *testing.T) {
	t.Run("successfully gets order", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		orderID := uuid.New()
		customerID := uuid.New()
		order := entity.NewOrder(customerID, 100.0, "pending")
		order.ID = orderID

		req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(orderID.String())

		mockRepo.On("FindByID", mock.Anything, orderID).Return(order, nil)

		err := h.GetByID(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 400 for invalid UUID", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		req := httptest.NewRequest(http.MethodGet, "/orders/invalid-uuid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid-uuid")

		err := h.GetByID(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		orderID := uuid.New()

		req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(orderID.String())

		mockRepo.On("FindByID", mock.Anything, orderID).Return(nil, errors.New("not found"))

		err := h.GetByID(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestOrderHandler_List(t *testing.T) {
	t.Run("successfully lists orders", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		orders := []entity.Order{
			*entity.NewOrder(uuid.New(), 100.0, "pending"),
			*entity.NewOrder(uuid.New(), 200.0, "confirmed"),
		}

		req := httptest.NewRequest(http.MethodGet, "/orders?offset=0&limit=10", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(orders, int64(2), nil)

		err := h.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 500 on repository error", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))

		err := h.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestOrderHandler_Update(t *testing.T) {
	t.Run("successfully updates order", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		orderID := uuid.New()
		customerID := uuid.New()
		reqBody := `{"customer_id":"` + customerID.String() + `","total":150.00,"status":"confirmed"}`

		req := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String(), strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(orderID.String())

		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

		err := h.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 400 for invalid UUID", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		req := httptest.NewRequest(http.MethodPut, "/orders/invalid-uuid", strings.NewReader(`{}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid-uuid")

		err := h.Update(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestOrderHandler_Delete(t *testing.T) {
	t.Run("successfully deletes order", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		orderID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/orders/"+orderID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(orderID.String())

		mockRepo.On("Delete", mock.Anything, orderID).Return(nil)

		err := h.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 400 for invalid UUID", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		req := httptest.NewRequest(http.MethodDelete, "/orders/invalid-uuid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid-uuid")

		err := h.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 500 on repository error", func(t *testing.T) {
		e, mockRepo := setupOrderHandlerTest()

		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		orderID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/orders/"+orderID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(orderID.String())

		mockRepo.On("Delete", mock.Anything, orderID).Return(errors.New("delete failed"))

		err := h.Delete(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockRepo.AssertExpectations(t)
	})
}

func TestOrderHandler_RegisterRoutes(t *testing.T) {
	t.Run("registers all routes", func(t *testing.T) {
		e := echo.New()
		mockRepo := new(MockOrderRepository)
		cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
		qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
		h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

		g := e.Group("/api/v1")
		h.RegisterRoutes(g)

		routes := e.Routes()
		require.NotEmpty(t, routes)

		// Verify expected routes exist
		routePaths := make(map[string]bool)
		for _, r := range routes {
			routePaths[r.Method+":"+r.Path] = true
		}

		assert.True(t, routePaths["POST:/api/v1/orders"])
		assert.True(t, routePaths["GET:/api/v1/orders"])
		assert.True(t, routePaths["GET:/api/v1/orders/:id"])
		assert.True(t, routePaths["PUT:/api/v1/orders/:id"])
		assert.True(t, routePaths["DELETE:/api/v1/orders/:id"])
	})
}

// =============================================================================
// Health Handler Tests
// =============================================================================

func TestHealthHandler(t *testing.T) {
	e := echo.New()

	t.Run("returns healthy status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Simulating health handler
		err := c.JSON(http.StatusOK, map[string]string{"status": "healthy"})

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "healthy", resp["status"])
	})
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkOrderHandler_Create(b *testing.B) {
	e := echo.New()
	e.Validator = validator.NewEchoValidator()
	mockRepo := new(MockOrderRepository)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Order")).Return(nil)

	cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
	qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
	h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

	customerID := uuid.New()
	reqBody := `{"customer_id":"` + customerID.String() + `","total":100.50,"status":"pending"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = h.Create(c)
	}
}

func BenchmarkOrderHandler_GetByID(b *testing.B) {
	e := echo.New()
	mockRepo := new(MockOrderRepository)

	orderID := uuid.New()
	order := entity.NewOrder(uuid.New(), 100.0, "pending")
	order.ID = orderID

	mockRepo.On("FindByID", mock.Anything, orderID).Return(order, nil)

	cmdHandler := apphandler.NewOrderCommandHandler(mockRepo)
	qryHandler := apphandler.NewOrderQueryHandler(mockRepo)
	h := httphandler.NewOrderHandler(cmdHandler, qryHandler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(orderID.String())
		_ = h.GetByID(c)
	}
}
