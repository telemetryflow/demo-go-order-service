// Package http provides HTTP routing.
package http

import (
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/telemetryflow/order-service/internal/infrastructure/http/handler"
	"github.com/telemetryflow/order-service/internal/infrastructure/http/middleware"
)

// setupRoutes configures all routes
func (s *Server) setupRoutes() {
	e := s.echo

	// Global middleware
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.RateLimit(s.config.RateLimit))

	// Health check
	healthHandler := handler.NewHealthHandler(s.db)
	e.GET("/health", healthHandler.Health)
	e.GET("/ready", healthHandler.Ready)

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

	// Swagger documentation
	// e.GET("/swagger/*", echoSwagger.WrapHandler)
}
