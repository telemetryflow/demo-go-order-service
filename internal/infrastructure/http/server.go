// Package http provides HTTP server implementation.
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

package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/telemetryflow/order-service/internal/infrastructure/config"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	echo   *echo.Echo
	config *config.Config
	db     *gorm.DB
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	server := &Server{
		echo:   e,
		config: cfg,
		db:     db,
	}

	// Setup routes
	server.setupRoutes()

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.echo.Start(":" + s.config.Server.Port)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

// Echo returns the underlying Echo instance
func (s *Server) Echo() *echo.Echo {
	return s.echo
}

// ServeHTTP implements http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.echo.ServeHTTP(w, r)
}
