// Package logs provides telemetry logging helpers.
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
	_ = telemetry.Client().LogInfo(ctx, message, attrs)
}

// Warn logs a warning-level message
func Warn(message string, attrs map[string]interface{}) {
	if !telemetry.IsEnabled() {
		log.Printf("[WARN] %s %v", message, attrs)
		return
	}
	ctx := context.Background()
	_ = telemetry.Client().LogWarn(ctx, message, attrs)
}

// Error logs an error-level message
func Error(message string, attrs map[string]interface{}) {
	if !telemetry.IsEnabled() {
		log.Printf("[ERROR] %s %v", message, attrs)
		return
	}
	ctx := context.Background()
	_ = telemetry.Client().LogError(ctx, message, attrs)
}

// Debug logs a debug-level message
func Debug(message string, attrs map[string]interface{}) {
	if !telemetry.IsEnabled() {
		log.Printf("[DEBUG] %s %v", message, attrs)
		return
	}
	ctx := context.Background()
	_ = telemetry.Client().Log(ctx, "debug", message, attrs)
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
