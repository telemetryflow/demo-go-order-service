// Package handler provides HTTP handlers.
package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HomeHandler handles home endpoint
type HomeHandler struct{}

// NewHomeHandler creates a new home handler
func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// Home returns service information
func (h *HomeHandler) Home(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"service": "Order-Service",
		"version": "1.1.2",
		"status":  "running",
		"message": "Welcome to Order Service API",
	})
}
