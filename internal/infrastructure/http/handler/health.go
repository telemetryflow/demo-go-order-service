// Package handler provides HTTP health check handler.
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

package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db        *gorm.DB
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
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
		sqlDB, err := h.db.DB()
		if err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			return c.JSON(http.StatusServiceUnavailable, HealthResponse{
				Status:    "unhealthy",
				Timestamp: time.Now(),
				Checks:    checks,
			})
		}
		if err := sqlDB.Ping(); err != nil {
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
