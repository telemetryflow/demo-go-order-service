// Package traces provides telemetry tracing helpers.
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

package traces

import (
	"context"

	"github.com/telemetryflow/order-service/telemetry"
)

// StartSpan starts a new trace span with server kind (for HTTP handlers)
func StartSpan(ctx context.Context, name string, attrs map[string]interface{}) (string, error) {
	if !telemetry.IsEnabled() {
		return "", nil
	}
	return telemetry.Client().StartSpan(ctx, name, "server", attrs)
}

// StartInternalSpan starts a new internal span (for internal operations)
func StartInternalSpan(ctx context.Context, name string, attrs map[string]interface{}) (string, error) {
	if !telemetry.IsEnabled() {
		return "", nil
	}
	return telemetry.Client().StartSpan(ctx, name, "internal", attrs)
}

// StartClientSpan starts a new client span (for outgoing requests)
func StartClientSpan(ctx context.Context, name string, attrs map[string]interface{}) (string, error) {
	if !telemetry.IsEnabled() {
		return "", nil
	}
	return telemetry.Client().StartSpan(ctx, name, "client", attrs)
}

// EndSpan ends an active span
func EndSpan(ctx context.Context, spanID string, spanErr error) error {
	if !telemetry.IsEnabled() || spanID == "" {
		return nil
	}
	return telemetry.Client().EndSpan(ctx, spanID, spanErr)
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, spanID, name string, attrs map[string]interface{}) error {
	if !telemetry.IsEnabled() || spanID == "" {
		return nil
	}
	return telemetry.Client().AddSpanEvent(ctx, spanID, name, attrs)
}

// SpanFunc wraps a function with tracing
func SpanFunc(ctx context.Context, name string, fn func() error) error {
	spanID, _ := StartSpan(ctx, name, nil)
	err := fn()
	_ = EndSpan(ctx, spanID, err)
	return err
}

// HTTPSpan creates a span for HTTP handlers
func HTTPSpan(ctx context.Context, method, path string) (string, error) {
	return StartSpan(ctx, "http."+method+"."+path, map[string]interface{}{
		"http.method": method,
		"http.path":   path,
	})
}

// DBSpan creates a span for database operations
func DBSpan(ctx context.Context, operation, table string) (string, error) {
	return StartSpan(ctx, "db."+operation+"."+table, map[string]interface{}{
		"db.operation": operation,
		"db.table":     table,
	})
}
