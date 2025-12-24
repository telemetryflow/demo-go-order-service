// Package logs provides telemetry logging helpers.
package logs

import (
	"context"
	"log"

	"github.com/telemetryflow/order-service/telemetry"
)

// Info logs an info-level message
func Info(message string, attrs map[string]interface{}) {
	if !telemetry.IsEnabled() {
		log.Printf("[INFO] %s %v", message, attrs)
		return
	}
	ctx := context.Background()
	telemetry.Client().LogInfo(ctx, message, attrs)
}

// Warn logs a warning-level message
func Warn(message string, attrs map[string]interface{}) {
	if !telemetry.IsEnabled() {
		log.Printf("[WARN] %s %v", message, attrs)
		return
	}
	ctx := context.Background()
	telemetry.Client().LogWarn(ctx, message, attrs)
}

// Error logs an error-level message
func Error(message string, attrs map[string]interface{}) {
	if !telemetry.IsEnabled() {
		log.Printf("[ERROR] %s %v", message, attrs)
		return
	}
	ctx := context.Background()
	telemetry.Client().LogError(ctx, message, attrs)
}

// Debug logs a debug-level message
func Debug(message string, attrs map[string]interface{}) {
	if !telemetry.IsEnabled() {
		log.Printf("[DEBUG] %s %v", message, attrs)
		return
	}
	ctx := context.Background()
	telemetry.Client().Log(ctx, "debug", message, attrs)
}

// WithError adds error to attributes
func WithError(err error) map[string]interface{} {
	if err == nil {
		return nil
	}
	return map[string]interface{}{
		"error": err.Error(),
	}
}

// Merge merges multiple attribute maps
func Merge(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
