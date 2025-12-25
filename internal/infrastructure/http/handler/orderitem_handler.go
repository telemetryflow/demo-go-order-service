// Package handler provides HTTP handlers for Orderitem.
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

// OrderitemHandler handles orderitem HTTP requests
type OrderitemHandler struct {
	commandHandler *handler.OrderitemCommandHandler
	queryHandler   *handler.OrderitemQueryHandler
}

// NewOrderitemHandler creates a new orderitem handler
func NewOrderitemHandler(
	cmdHandler *handler.OrderitemCommandHandler,
	qryHandler *handler.OrderitemQueryHandler,
) *OrderitemHandler {
	return &OrderitemHandler{
		commandHandler: cmdHandler,
		queryHandler:   qryHandler,
	}
}

// RegisterRoutes registers orderitem routes
func (h *OrderitemHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/orderitems", h.Create)
	g.GET("/orderitems", h.List)
	g.GET("/orderitems/:id", h.GetByID)
	g.PUT("/orderitems/:id", h.Update)
	g.DELETE("/orderitems/:id", h.Delete)
}

// Create handles POST /orderitems
func (h *OrderitemHandler) Create(c echo.Context) error {
	var req dto.CreateOrderitemRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	cmd := &command.CreateOrderitemCommand{
		OrderId:   req.OrderId,
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
		Price:     req.Price,
	}

	if err := h.commandHandler.HandleOrderitemCreate(c.Request().Context(), cmd); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, nil, "Orderitem created successfully")
}

// List handles GET /orderitems
func (h *OrderitemHandler) List(c echo.Context) error {
	var q query.GetAllorderitemsQuery
	if err := c.Bind(&q); err != nil {
		return response.BadRequest(c, "Invalid query parameters")
	}
	_ = q.Validate()

	result, err := h.queryHandler.HandleOrderitemGetAll(c.Request().Context(), &q)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, result, "")
}

// GetByID handles GET /orderitems/:id
func (h *OrderitemHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	q := &query.GetOrderitemByIDQuery{ID: id}
	result, err := h.queryHandler.HandleOrderitemGetByID(c.Request().Context(), q)
	if err != nil {
		return response.NotFound(c, "Orderitem not found")
	}

	return response.Success(c, result, "")
}

// Update handles PUT /orderitems/:id
func (h *OrderitemHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	var req dto.UpdateOrderitemRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	cmd := &command.UpdateOrderitemCommand{
		ID:        id,
		OrderId:   req.OrderId,
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
		Price:     req.Price,
	}

	if err := h.commandHandler.HandleOrderitemUpdate(c.Request().Context(), cmd); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, nil, "Orderitem updated successfully")
}

// Delete handles DELETE /orderitems/:id
func (h *OrderitemHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	cmd := &command.DeleteOrderitemCommand{ID: id}
	if err := h.commandHandler.HandleOrderitemDelete(c.Request().Context(), cmd); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}
