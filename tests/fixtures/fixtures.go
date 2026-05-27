// Package fixtures provides test fixtures and sample data.
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

package fixtures

import (
	"os"

	"github.com/google/uuid"
	"github.com/telemetryflow/order-service/internal/domain"
)

// SetEnvVars sets environment variables for testing and returns a cleanup function
func SetEnvVars(vars map[string]string) func() {
	originalVars := make(map[string]string)

	// Store original values
	for key, value := range vars {
		originalVars[key] = os.Getenv(key)
		_ = os.Setenv(key, value)
	}

	// Return cleanup function
	return func() {
		for key, originalValue := range originalVars {
			if originalValue == "" {
				_ = os.Unsetenv(key)
			} else {
				_ = os.Setenv(key, originalValue)
			}
		}
	}
}

// GetSampleEntity returns a sample entity for testing
func GetSampleEntity() *domain.Entity {
	return domain.NewEntity(uuid.New(), "Test Entity")
}

// GetSampleEntities returns a slice of sample entities for testing
func GetSampleEntities(count int) []interface{} {
	entities := make([]interface{}, count)
	for i := 0; i < count; i++ {
		entities[i] = domain.NewEntity(uuid.New(), "Test Entity")
	}
	return entities
}

// GetSampleEntityData returns sample entity data for HTTP request body
func GetSampleEntityData() map[string]interface{} {
	return map[string]interface{}{
		"name": "Test Entity",
	}
}
