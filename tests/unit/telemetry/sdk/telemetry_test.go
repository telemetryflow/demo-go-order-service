// Package sdk_test provides unit tests for TelemetryFlow SDK integration.
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

package sdk_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telemetryflow/order-service/telemetry"
)

// =============================================================================
// Telemetry Tests
//
// Tests for TelemetryFlow SDK initialization and lifecycle management.
// =============================================================================

// TestTelemetry_Init_WithoutCredentials verifies graceful degradation without API keys.
func TestTelemetry_Init_WithoutCredentials(t *testing.T) {
	// Clear any existing credentials (t.Setenv to empty clears and auto-restores)
	t.Setenv("TELEMETRYFLOW_API_KEY_ID", "")
	t.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "")

	t.Run("init returns nil without credentials", func(t *testing.T) {
		err := telemetry.Init()
		assert.NoError(t, err)
	})

	t.Run("client is nil without credentials", func(t *testing.T) {
		// After init without credentials
		client := telemetry.Client()
		assert.Nil(t, client)
	})

	t.Run("IsEnabled returns false without credentials", func(t *testing.T) {
		assert.False(t, telemetry.IsEnabled())
	})
}

func TestTelemetry_Shutdown(t *testing.T) {
	t.Run("shutdown does not panic when client is nil", func(t *testing.T) {
		// Clear credentials (t.Setenv to empty clears and auto-restores)
		t.Setenv("TELEMETRYFLOW_API_KEY_ID", "")
		t.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "")

		// Init without credentials (client will be nil)
		_ = telemetry.Init()

		// Shutdown should not panic
		assert.NotPanics(t, func() {
			telemetry.Shutdown()
		})
	})
}

func TestTelemetry_Client(t *testing.T) {
	t.Run("returns nil when not initialized", func(t *testing.T) {
		t.Setenv("TELEMETRYFLOW_API_KEY_ID", "")
		t.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "")

		_ = telemetry.Init()

		client := telemetry.Client()
		assert.Nil(t, client)
	})
}

func TestTelemetry_IsEnabled(t *testing.T) {
	tests := []struct {
		name      string
		keyID     string
		keySecret string
		expected  bool
	}{
		{
			name:      "disabled when no credentials",
			keyID:     "",
			keySecret: "",
			expected:  false,
		},
		{
			name:      "disabled when only key ID",
			keyID:     "tfk_test",
			keySecret: "",
			expected:  false,
		},
		{
			name:      "disabled when only key secret",
			keyID:     "",
			keySecret: "tfs_test",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Setenv handles cleanup automatically after the test
			t.Setenv("TELEMETRYFLOW_API_KEY_ID", tt.keyID)
			t.Setenv("TELEMETRYFLOW_API_KEY_SECRET", tt.keySecret)

			_ = telemetry.Init()
			assert.Equal(t, tt.expected, telemetry.IsEnabled())
		})
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestTelemetry_EdgeCases(t *testing.T) {
	t.Run("multiple init calls do not panic", func(t *testing.T) {
		t.Setenv("TELEMETRYFLOW_API_KEY_ID", "")
		t.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "")

		assert.NotPanics(t, func() {
			_ = telemetry.Init()
			_ = telemetry.Init()
			_ = telemetry.Init()
		})
	})

	t.Run("multiple shutdown calls do not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			telemetry.Shutdown()
			telemetry.Shutdown()
			telemetry.Shutdown()
		})
	})

	t.Run("init then shutdown cycle", func(t *testing.T) {
		t.Setenv("TELEMETRYFLOW_API_KEY_ID", "")
		t.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "")

		assert.NotPanics(t, func() {
			for i := 0; i < 3; i++ {
				_ = telemetry.Init()
				telemetry.Shutdown()
			}
		})
	})
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkTelemetry_IsEnabled(b *testing.B) {
	b.Setenv("TELEMETRYFLOW_API_KEY_ID", "")
	b.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "")
	_ = telemetry.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = telemetry.IsEnabled()
	}
}

func BenchmarkTelemetry_Client(b *testing.B) {
	b.Setenv("TELEMETRYFLOW_API_KEY_ID", "")
	b.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "")
	_ = telemetry.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = telemetry.Client()
	}
}
