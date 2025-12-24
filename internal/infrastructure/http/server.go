// Package http provides HTTP server implementation.
package http

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/telemetryflow/order-service/internal/infrastructure/config"
)

// Server represents the HTTP server
type Server struct {
	echo   *echo.Echo
	config *config.Config
	db     *sql.DB
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, db *sql.DB) *Server {
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
