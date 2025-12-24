// Package telemetry provides TelemetryFlow SDK integration.
package telemetry

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

var client *telemetryflow.Client

// Init initializes the TelemetryFlow SDK
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

	client, err = telemetryflow.NewBuilder().
		WithAutoConfiguration().
		WithInsecure(insecure).
		WithSignals(true, true, true).
		Build()
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		return err
	}

	log.Println("TelemetryFlow SDK initialized successfully")
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
