// Package http provides HTTP route registration.
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
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/telemetryflow/order-service/internal/infrastructure/http/handler"
	"github.com/telemetryflow/order-service/internal/infrastructure/http/middleware"
)

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	e := s.echo

	// Global middleware
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())

	// OpenTelemetry auto-instrumentation for HTTP
	e.Use(otelecho.Middleware(s.config.Telemetry.ServiceName))

	// Set span status based on HTTP response code
	e.Use(spanStatusMiddleware())

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimit(s.config.RateLimit))

	// Health check
	healthHandler := handler.NewHealthHandler(s.db)
	e.GET("/health", healthHandler.Health)
	e.GET("/ready", healthHandler.Ready)

	// Home endpoint
	homeHandler := handler.NewHomeHandler()
	e.GET("/", homeHandler.Home)

	// Swagger documentation
	swaggerHandler := handler.NewSwaggerHandler("Order Service API")
	swaggerHandler.RegisterRoutes(e)

	// API v1 routes
	v1 := e.Group("/api/v1")
	{
		// Public routes
		// v1.POST("/auth/login", authHandler.Login)
		// v1.POST("/auth/register", authHandler.Register)

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.Auth(s.config.JWT))
		{
			// Add protected routes here
		}
	}

}

// spanStatusMiddleware sets the span status based on HTTP response status code
func spanStatusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Process request first
			err := next(c)

			// Get the current span from context and set status
			span := trace.SpanFromContext(c.Request().Context())
			if span.IsRecording() {
				status := c.Response().Status
				if status >= 500 {
					span.SetStatus(codes.Error, "Server error")
				} else if status >= 400 {
					span.SetStatus(codes.Error, "Client error")
				} else {
					span.SetStatus(codes.Ok, "Success")
				}
			}

			return err
		}
	}
}
