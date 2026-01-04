// Package main is the entry point for Order-Service API.
//
// Order-Service - RESTful API with DDD + CQRS Pattern
// Copyright (c) 2024-2026 Order-Service. All rights reserved.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/telemetryflow/order-service/internal/infrastructure/config"
	"github.com/telemetryflow/order-service/internal/infrastructure/http"
	"github.com/telemetryflow/order-service/internal/infrastructure/persistence"
	"github.com/telemetryflow/order-service/telemetry"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	// Initialize TelemetryFlow
	if err := telemetry.Init(); err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer telemetry.Shutdown()

	// Initialize database
	db, err := persistence.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Failed to get underlying sql.DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	// Create HTTP server
	server := http.NewServer(cfg, db)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	log.Printf("Order-Service API v1.1.2 started on port %s", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
