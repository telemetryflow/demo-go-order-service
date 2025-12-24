// Package metrics provides telemetry metrics helpers.
package metrics

import (
	"context"

	"github.com/telemetryflow/order-service/telemetry"
)

// IncrementCounter increments a counter metric
func IncrementCounter(name string, value int64, labels map[string]interface{}) {
	if !telemetry.IsEnabled() {
		return
	}
	ctx := context.Background()
	telemetry.Client().IncrementCounter(ctx, name, value, labels)
}

// RecordGauge records a gauge metric
func RecordGauge(name string, value float64, labels map[string]interface{}) {
	if !telemetry.IsEnabled() {
		return
	}
	ctx := context.Background()
	telemetry.Client().RecordGauge(ctx, name, value, labels)
}

// RecordHistogram records a histogram measurement
func RecordHistogram(name string, value float64, unit string, labels map[string]interface{}) {
	if !telemetry.IsEnabled() {
		return
	}
	ctx := context.Background()
	telemetry.Client().RecordHistogram(ctx, name, value, unit, labels)
}

// HTTP Metrics

// RecordHTTPRequest records an HTTP request metric
func RecordHTTPRequest(method, path string, statusCode int, duration float64) {
	RecordHistogram("http.request.duration", duration, "s", map[string]interface{}{
		"method": method,
		"path":   path,
		"status": statusCode,
	})
	IncrementCounter("http.requests.total", 1, map[string]interface{}{
		"method": method,
		"path":   path,
		"status": statusCode,
	})
}

// Database Metrics

// RecordDBQuery records a database query metric
func RecordDBQuery(operation, table string, duration float64, success bool) {
	RecordHistogram("db.query.duration", duration, "s", map[string]interface{}{
		"operation": operation,
		"table":     table,
		"success":   success,
	})
	IncrementCounter("db.queries.total", 1, map[string]interface{}{
		"operation": operation,
		"table":     table,
		"success":   success,
	})
}

// Business Metrics

// RecordEntityCreated records an entity creation
func RecordEntityCreated(entityType string) {
	IncrementCounter("entity.created.total", 1, map[string]interface{}{
		"type": entityType,
	})
}

// RecordEntityUpdated records an entity update
func RecordEntityUpdated(entityType string) {
	IncrementCounter("entity.updated.total", 1, map[string]interface{}{
		"type": entityType,
	})
}

// RecordEntityDeleted records an entity deletion
func RecordEntityDeleted(entityType string) {
	IncrementCounter("entity.deleted.total", 1, map[string]interface{}{
		"type": entityType,
	})
}
