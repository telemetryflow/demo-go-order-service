// Package validator_test provides unit tests for validators.
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

package validator_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/telemetryflow/order-service/pkg/validator"
)

// =============================================================================
// Test Structs
//
// Sample request structures with validation tags for testing.
// =============================================================================

type TestRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"gte=0,lte=150"`
	Website  string `json:"website" validate:"omitempty,url"`
	Role     string `json:"role" validate:"required,oneof=admin user guest"`
	Password string `json:"password" validate:"required,min=8"`
	UUID     string `json:"uuid" validate:"omitempty,uuid"`
}

type SimpleRequest struct {
	Field string `json:"field" validate:"required"`
}

// =============================================================================
// Validator Tests
// =============================================================================

func TestNew(t *testing.T) {
	t.Run("creates new validator", func(t *testing.T) {
		v := validator.New()
		require.NotNil(t, v)
	})
}

func TestValidator_Validate_Success(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name: "valid complete request",
			input: &TestRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				Website:  "https://example.com",
				Role:     "admin",
				Password: "password123",
				UUID:     uuid.New().String(),
			},
		},
		{
			name: "valid minimal request",
			input: &TestRequest{
				Name:     "Jo",
				Email:    "j@e.co",
				Age:      0,
				Role:     "user",
				Password: "12345678",
			},
		},
		{
			name: "valid with max values",
			input: &TestRequest{
				Name:     "A very long name that is still valid", // Max 50
				Email:    "long.email.address@subdomain.example.com",
				Age:      150,
				Role:     "guest",
				Password: "averylongpassword123456789",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.input)
			assert.NoError(t, err)
		})
	}
}

func TestValidator_Validate_Required(t *testing.T) {
	v := validator.New()

	t.Run("missing required field", func(t *testing.T) {
		input := &SimpleRequest{Field: ""}

		err := v.Validate(input)

		require.Error(t, err)
		validationErr, ok := err.(*validator.ValidationError)
		require.True(t, ok)
		assert.Contains(t, validationErr.Errors["field"], "required")
	})
}

func TestValidator_Validate_Email(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		email   string
		isValid bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@mail.example.com", true},
		{"invalid email no @", "testexample.com", false},
		{"invalid email no domain", "test@", false},
		{"invalid email spaces", "test @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &TestRequest{
				Name:     "Test",
				Email:    tt.email,
				Role:     "admin",
				Password: "password123",
			}

			err := v.Validate(input)

			if tt.isValid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				validationErr, ok := err.(*validator.ValidationError)
				require.True(t, ok)
				assert.Contains(t, validationErr.Errors["email"], "email")
			}
		})
	}
}

func TestValidator_Validate_MinMax(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		field   string
		value   string
		isValid bool
	}{
		{"name at min length", "name", "Jo", true},
		{"name below min length", "name", "J", false},
		{"name at max length", "name", "ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUV", true}, // 48 chars
		{"password at min length", "password", "12345678", true},
		{"password below min length", "password", "1234567", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &TestRequest{
				Name:     "Test Name",
				Email:    "test@example.com",
				Role:     "admin",
				Password: "password123",
			}

			switch tt.field {
			case "name":
				input.Name = tt.value
			case "password":
				input.Password = tt.value
			}

			err := v.Validate(input)

			if tt.isValid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				validationErr, ok := err.(*validator.ValidationError)
				require.True(t, ok)
				assert.NotEmpty(t, validationErr.Errors[tt.field])
			}
		})
	}
}

func TestValidator_Validate_OneOf(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		role    string
		isValid bool
	}{
		{"valid role admin", "admin", true},
		{"valid role user", "user", true},
		{"valid role guest", "guest", true},
		{"invalid role", "superadmin", false},
		{"empty role", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &TestRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Role:     tt.role,
				Password: "password123",
			}

			err := v.Validate(input)

			if tt.isValid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestValidator_Validate_GteLte(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		age     int
		isValid bool
	}{
		{"age at min (0)", 0, true},
		{"age at max (150)", 150, true},
		{"age in range", 30, true},
		{"age below min", -1, false},
		{"age above max", 151, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &TestRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Age:      tt.age,
				Role:     "admin",
				Password: "password123",
			}

			err := v.Validate(input)

			if tt.isValid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				validationErr, ok := err.(*validator.ValidationError)
				require.True(t, ok)
				assert.NotEmpty(t, validationErr.Errors["age"])
			}
		})
	}
}

func TestValidator_Validate_URL(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		url     string
		isValid bool
	}{
		{"valid https url", "https://example.com", true},
		{"valid http url", "http://example.com", true},
		{"valid url with path", "https://example.com/path", true},
		{"empty url (optional)", "", true},
		{"invalid url", "not-a-url", false},
		{"invalid url missing protocol", "example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &TestRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Website:  tt.url,
				Role:     "admin",
				Password: "password123",
			}

			err := v.Validate(input)

			if tt.isValid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				validationErr, ok := err.(*validator.ValidationError)
				require.True(t, ok)
				assert.NotEmpty(t, validationErr.Errors["website"])
			}
		})
	}
}

func TestValidator_Validate_UUID(t *testing.T) {
	v := validator.New()

	tests := []struct {
		name    string
		uuid    string
		isValid bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid UUID lowercase", "550e8400-e29b-41d4-a716-446655440000", true},
		{"empty UUID (optional)", "", true},
		{"invalid UUID format", "not-a-uuid", false},
		{"invalid UUID too short", "550e8400-e29b-41d4", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &TestRequest{
				Name:     "Test",
				Email:    "test@example.com",
				Role:     "admin",
				Password: "password123",
				UUID:     tt.uuid,
			}

			err := v.Validate(input)

			if tt.isValid {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				validationErr, ok := err.(*validator.ValidationError)
				require.True(t, ok)
				assert.NotEmpty(t, validationErr.Errors["uuid"])
			}
		})
	}
}

// =============================================================================
// ValidationError Tests
// =============================================================================

func TestValidationError_Error(t *testing.T) {
	t.Run("returns message", func(t *testing.T) {
		err := &validator.ValidationError{
			Message: "Validation failed",
			Errors:  map[string]string{"field": "required"},
		}

		assert.Equal(t, "Validation failed", err.Error())
	})
}

func TestValidationError_HTTPError(t *testing.T) {
	t.Run("returns HTTP 400 error", func(t *testing.T) {
		validationErr := &validator.ValidationError{
			Message: "Validation failed",
			Errors:  map[string]string{"field": "required"},
		}

		httpErr := validationErr.HTTPError()

		require.NotNil(t, httpErr)
		assert.Equal(t, 400, httpErr.Code)
	})
}

// =============================================================================
// EchoValidator Tests
// =============================================================================

func TestNewEchoValidator(t *testing.T) {
	t.Run("creates new Echo validator", func(t *testing.T) {
		v := validator.NewEchoValidator()
		require.NotNil(t, v)
	})
}

func TestEchoValidator_Validate(t *testing.T) {
	v := validator.NewEchoValidator()

	t.Run("validates successfully", func(t *testing.T) {
		input := &SimpleRequest{Field: "value"}
		err := v.Validate(input)
		assert.NoError(t, err)
	})

	t.Run("returns error for invalid input", func(t *testing.T) {
		input := &SimpleRequest{Field: ""}
		err := v.Validate(input)
		assert.Error(t, err)
	})
}

// =============================================================================
// Multiple Errors Tests
// =============================================================================

func TestValidator_Validate_MultipleErrors(t *testing.T) {
	v := validator.New()

	t.Run("returns multiple validation errors", func(t *testing.T) {
		input := &TestRequest{
			Name:     "",        // required, min=2
			Email:    "invalid", // email format
			Age:      200,       // lte=150
			Role:     "invalid", // oneof
			Password: "short",   // min=8
		}

		err := v.Validate(input)

		require.Error(t, err)
		validationErr, ok := err.(*validator.ValidationError)
		require.True(t, ok)
		assert.Equal(t, "Validation failed", validationErr.Message)

		// Should have multiple field errors
		assert.NotEmpty(t, validationErr.Errors)
	})
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkValidator_Validate_Valid(b *testing.B) {
	v := validator.New()
	input := &TestRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		Role:     "admin",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Validate(input)
	}
}

func BenchmarkValidator_Validate_Invalid(b *testing.B) {
	v := validator.New()
	input := &TestRequest{
		Name:     "",
		Email:    "invalid",
		Role:     "invalid",
		Password: "short",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Validate(input)
	}
}

func BenchmarkEchoValidator_Validate(b *testing.B) {
	v := validator.NewEchoValidator()
	input := &SimpleRequest{Field: "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Validate(input)
	}
}
