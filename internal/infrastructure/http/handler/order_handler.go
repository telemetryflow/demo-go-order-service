// Package handler provides HTTP handlers for Order endpoints.
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

package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/telemetryflow/order-service/internal/application/command"
	"github.com/telemetryflow/order-service/internal/application/dto"
	"github.com/telemetryflow/order-service/internal/application/handler"
	"github.com/telemetryflow/order-service/internal/application/query"
	"github.com/telemetryflow/order-service/pkg/response"
)

// OrderHandler handles order HTTP requests
type OrderHandler struct {
	commandHandler *handler.OrderCommandHandler
	queryHandler   *handler.OrderQueryHandler
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(
	cmdHandler *handler.OrderCommandHandler,
	qryHandler *handler.OrderQueryHandler,
) *OrderHandler {
	return &OrderHandler{
		commandHandler: cmdHandler,
		queryHandler:   qryHandler,
	}
}

// RegisterRoutes registers order routes
func (h *OrderHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/orders", h.Create)
	g.GET("/orders", h.List)
	g.GET("/orders/:id", h.GetByID)
	g.PUT("/orders/:id", h.Update)
	g.DELETE("/orders/:id", h.Delete)
}

// Create handles POST /orders
func (h *OrderHandler) Create(c echo.Context) error {
	var req dto.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	cmd := &command.CreateOrderCommand{
		CustomerID: req.CustomerID,
		Total:      req.Total,
		Status:     req.Status,
	}

	if err := h.commandHandler.HandleOrderCreate(c.Request().Context(), cmd); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, nil, "Order created successfully")
}

// List handles GET /orders
func (h *OrderHandler) List(c echo.Context) error {
	var q query.GetAllOrdersQuery
	if err := c.Bind(&q); err != nil {
		return response.BadRequest(c, "Invalid query parameters")
	}
	_ = q.Validate()

	result, err := h.queryHandler.HandleOrderGetAll(c.Request().Context(), &q)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, result, "")
}

// GetByID handles GET /orders/:id
func (h *OrderHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	q := &query.GetOrderByIDQuery{ID: id}
	result, err := h.queryHandler.HandleOrderGetByID(c.Request().Context(), q)
	if err != nil {
		return response.NotFound(c, "Order not found")
	}

	return response.Success(c, result, "")
}

// Update handles PUT /orders/:id
func (h *OrderHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	var req dto.UpdateOrderRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	cmd := &command.UpdateOrderCommand{
		ID:         id,
		CustomerID: req.CustomerID,
		Total:      req.Total,
		Status:     req.Status,
	}

	if err := h.commandHandler.HandleOrderUpdate(c.Request().Context(), cmd); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, nil, "Order updated successfully")
}

// Delete handles DELETE /orders/:id
func (h *OrderHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	cmd := &command.DeleteOrderCommand{ID: id}
	if err := h.commandHandler.HandleOrderDelete(c.Request().Context(), cmd); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}
