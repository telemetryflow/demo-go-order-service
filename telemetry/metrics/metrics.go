// Package metrics provides telemetry metrics helpers.
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
	_ = telemetry.Client().IncrementCounter(ctx, name, value, labels)
}

// RecordGauge records a gauge metric
func RecordGauge(name string, value float64, labels map[string]interface{}) {
	if !telemetry.IsEnabled() {
		return
	}
	ctx := context.Background()
	_ = telemetry.Client().RecordGauge(ctx, name, value, labels)
}

// RecordHistogram records a histogram measurement
func RecordHistogram(name string, value float64, unit string, labels map[string]interface{}) {
	if !telemetry.IsEnabled() {
		return
	}
	ctx := context.Background()
	_ = telemetry.Client().RecordHistogram(ctx, name, value, unit, labels)
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
