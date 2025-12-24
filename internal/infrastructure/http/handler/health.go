// Package handler provides HTTP handlers.
package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db        *sql.DB
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// Health handles the health check endpoint
func (h *HealthHandler) Health(c echo.Context) error {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
	}
	return c.JSON(http.StatusOK, response)
}

// Ready handles the readiness check endpoint
func (h *HealthHandler) Ready(c echo.Context) error {
	checks := make(map[string]string)

	// Check database connection
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			return c.JSON(http.StatusServiceUnavailable, HealthResponse{
				Status:    "unhealthy",
				Timestamp: time.Now(),
				Checks:    checks,
			})
		}
		checks["database"] = "healthy"
	}

	return c.JSON(http.StatusOK, HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
		Uptime:    time.Since(h.startTime).String(),
		Checks:    checks,
	})
}
