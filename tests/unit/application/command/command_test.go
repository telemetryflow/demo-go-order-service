// Package command_test provides unit tests for CQRS commands.
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

package command_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/telemetryflow/order-service/internal/application/command"
)

// =============================================================================
// CreateOrderCommand Tests
//
// Tests for the CreateOrderCommand which handles new order creation requests.
// This command converts request data into a domain entity for persistence.
// =============================================================================

// TestCreateOrderCommand_Validate verifies validation rules for order creation.
func TestCreateOrderCommand_Validate(t *testing.T) {
	t.Run("valid command returns nil", func(t *testing.T) {
		cmd := &command.CreateOrderCommand{
			CustomerID: uuid.New(),
			Total:      100.0,
			Status:     "pending",
		}

		err := cmd.Validate()
		assert.NoError(t, err)
	})
}

func TestCreateOrderCommand_ToEntity(t *testing.T) {
	t.Run("converts to entity correctly", func(t *testing.T) {
		customerID := uuid.New()
		cmd := &command.CreateOrderCommand{
			CustomerID: customerID,
			Total:      150.50,
			Status:     "confirmed",
		}

		entity := cmd.ToEntity()

		require.NotNil(t, entity)
		assert.NotEqual(t, uuid.Nil, entity.ID)
		assert.Equal(t, customerID, entity.CustomerID)
		assert.Equal(t, 150.50, entity.Total)
		assert.Equal(t, "confirmed", entity.Status)
		assert.False(t, entity.CreatedAt.IsZero())
		assert.False(t, entity.UpdatedAt.IsZero())
	})

	t.Run("creates unique entities", func(t *testing.T) {
		cmd := &command.CreateOrderCommand{
			CustomerID: uuid.New(),
			Total:      100.0,
			Status:     "pending",
		}

		entity1 := cmd.ToEntity()
		entity2 := cmd.ToEntity()

		assert.NotEqual(t, entity1.ID, entity2.ID)
	})
}

// =============================================================================
// UpdateOrderCommand Tests
//
// Tests for the UpdateOrderCommand which handles order modification requests.
// This command requires a valid ID and validates all update fields.
// =============================================================================

// TestUpdateOrderCommand_Validate verifies validation rules for order updates.
func TestUpdateOrderCommand_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *command.UpdateOrderCommand
		expectError bool
	}{
		{
			name: "valid command returns nil",
			cmd: &command.UpdateOrderCommand{
				ID:         uuid.New(),
				CustomerID: uuid.New(),
				Total:      100.0,
				Status:     "pending",
			},
			expectError: false,
		},
		{
			name: "nil ID returns error",
			cmd: &command.UpdateOrderCommand{
				ID:         uuid.Nil,
				CustomerID: uuid.New(),
				Total:      100.0,
				Status:     "pending",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, command.ErrInvalidID, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateOrderCommand_ToEntity(t *testing.T) {
	t.Run("converts to entity with correct ID", func(t *testing.T) {
		id := uuid.New()
		customerID := uuid.New()
		cmd := &command.UpdateOrderCommand{
			ID:         id,
			CustomerID: customerID,
			Total:      200.0,
			Status:     "shipped",
		}

		entity := cmd.ToEntity()

		require.NotNil(t, entity)
		assert.Equal(t, id, entity.ID)
		assert.Equal(t, customerID, entity.CustomerID)
		assert.Equal(t, 200.0, entity.Total)
		assert.Equal(t, "shipped", entity.Status)
	})
}

// =============================================================================
// DeleteOrderCommand Tests
//
// Tests for the DeleteOrderCommand which handles order deletion requests.
// This command only requires a valid order ID for the soft delete operation.
// =============================================================================

// TestDeleteOrderCommand_Validate verifies validation rules for order deletion.
func TestDeleteOrderCommand_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *command.DeleteOrderCommand
		expectError bool
	}{
		{
			name: "valid ID returns nil",
			cmd: &command.DeleteOrderCommand{
				ID: uuid.New(),
			},
			expectError: false,
		},
		{
			name: "nil ID returns error",
			cmd: &command.DeleteOrderCommand{
				ID: uuid.Nil,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, command.ErrInvalidID, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =============================================================================
// Edge Cases
//
// Tests for unusual inputs and boundary conditions in command handling.
// =============================================================================

// TestCommand_EdgeCases verifies command behavior with edge case inputs.
func TestCommand_EdgeCases(t *testing.T) {
	t.Run("create command with zero total", func(t *testing.T) {
		cmd := &command.CreateOrderCommand{
			CustomerID: uuid.New(),
			Total:      0,
			Status:     "pending",
		}

		entity := cmd.ToEntity()
		assert.Equal(t, float64(0), entity.Total)
	})

	t.Run("create command with negative total", func(t *testing.T) {
		cmd := &command.CreateOrderCommand{
			CustomerID: uuid.New(),
			Total:      -50.0,
			Status:     "refunded",
		}

		entity := cmd.ToEntity()
		assert.Equal(t, -50.0, entity.Total)
	})

	t.Run("create command with empty status", func(t *testing.T) {
		cmd := &command.CreateOrderCommand{
			CustomerID: uuid.New(),
			Total:      100.0,
			Status:     "",
		}

		entity := cmd.ToEntity()
		assert.Equal(t, "", entity.Status)
	})

	t.Run("update command preserves provided ID", func(t *testing.T) {
		specificID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
		cmd := &command.UpdateOrderCommand{
			ID:         specificID,
			CustomerID: uuid.New(),
			Total:      100.0,
			Status:     "pending",
		}

		entity := cmd.ToEntity()
		assert.Equal(t, specificID, entity.ID)
	})
}

// =============================================================================
// Benchmark Tests
//
// Performance benchmarks for command operations.
// Run with: go test -bench=. -benchmem
// =============================================================================

// BenchmarkCreateOrderCommand_ToEntity measures entity conversion performance.
func BenchmarkCreateOrderCommand_ToEntity(b *testing.B) {
	cmd := &command.CreateOrderCommand{
		CustomerID: uuid.New(),
		Total:      100.0,
		Status:     "pending",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cmd.ToEntity()
	}
}

func BenchmarkUpdateOrderCommand_Validate(b *testing.B) {
	cmd := &command.UpdateOrderCommand{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
		Total:      100.0,
		Status:     "pending",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cmd.Validate()
	}
}

func BenchmarkDeleteOrderCommand_Validate(b *testing.B) {
	cmd := &command.DeleteOrderCommand{
		ID: uuid.New(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cmd.Validate()
	}
}
