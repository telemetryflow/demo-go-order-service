// Package telemetry provides TelemetryFlow SDK integration for Order-Service.
//
// TelemetryFlow Go SDK v1.1.2 - Compatible with TFO-Collector v1.1.2 (OCB-native)
package telemetry

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

var client *telemetryflow.Client

// Init initializes the TelemetryFlow SDK with TFO v2 API support.
// Configuration is loaded from environment variables:
//   - TELEMETRYFLOW_API_KEY_ID, TELEMETRYFLOW_API_KEY_SECRET (required)
//   - TELEMETRYFLOW_ENDPOINT (default: api.telemetryflow.id:4317)
//   - TELEMETRYFLOW_SERVICE_NAME, TELEMETRYFLOW_SERVICE_VERSION
//   - TELEMETRYFLOW_USE_V2_API, TELEMETRYFLOW_V2_ONLY
//   - TELEMETRYFLOW_COLLECTOR_NAME, TELEMETRYFLOW_DATACENTER
func Init() error {
	var err error

	// Check if credentials are provided
	keyID := os.Getenv("TELEMETRYFLOW_API_KEY_ID")
	keySecret := os.Getenv("TELEMETRYFLOW_API_KEY_SECRET")

	if keyID == "" || keySecret == "" {
		log.Println("TelemetryFlow credentials not found, telemetry disabled")
		return nil
	}

	// Check if insecure mode is enabled (for local development)
	insecure := os.Getenv("TELEMETRYFLOW_INSECURE") == "true"

	// Check for v2-only mode
	v2Only := os.Getenv("TELEMETRYFLOW_V2_ONLY") == "true"

	// Build client with TFO v2 API and collector identity support
	builder := telemetryflow.NewBuilder().
		WithAutoConfiguration().
		WithInsecure(insecure).
		WithSignals(true, true, true).
		WithExemplars(true)

	// Enable v2-only mode if configured
	if v2Only {
		builder = builder.WithV2Only()
	}

	client, err = builder.Build()
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		return err
	}

	log.Println("TelemetryFlow SDK v1.1.2 initialized successfully (TFO v2 API enabled)")
	return nil
}

// InitWithV2Only initializes the TelemetryFlow SDK in v2-only mode.
// This mode uses only TFO Platform v2 endpoints for maximum compatibility
// with TFO-Collector v1.1.2 (OCB-native).
func InitWithV2Only() error {
	var err error

	keyID := os.Getenv("TELEMETRYFLOW_API_KEY_ID")
	keySecret := os.Getenv("TELEMETRYFLOW_API_KEY_SECRET")

	if keyID == "" || keySecret == "" {
		log.Println("TelemetryFlow credentials not found, telemetry disabled")
		return nil
	}

	insecure := os.Getenv("TELEMETRYFLOW_INSECURE") == "true"

	client, err = telemetryflow.NewBuilder().
		WithAutoConfiguration().
		WithInsecure(insecure).
		WithV2Only(). // Enable v2-only mode
		WithSignals(true, true, true).
		WithExemplars(true).
		Build()
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		return err
	}

	log.Println("TelemetryFlow SDK v1.1.2 initialized in v2-only mode")
	return nil
}

// Shutdown gracefully shuts down the SDK
func Shutdown() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down TelemetryFlow: %v", err)
		}
	}
}

// Client returns the global TelemetryFlow client
func Client() *telemetryflow.Client {
	return client
}

// IsEnabled returns true if telemetry is enabled
func IsEnabled() bool {
	return client != nil
}
