// Package response_test provides unit tests for response helpers.
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

package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/telemetryflow/order-service/pkg/response"
)

// =============================================================================
// Helper Functions
//
// Test utilities for parsing and validating response bodies.
// =============================================================================

// parseResponse is a helper to unmarshal JSON response bodies.
func parseResponse(t *testing.T, rec *httptest.ResponseRecorder) response.Response {
	var resp response.Response
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)
	return resp
}

// =============================================================================
// Success Response Tests
// =============================================================================

func TestSuccess(t *testing.T) {
	e := echo.New()

	t.Run("returns 200 with data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		data := map[string]string{"key": "value"}
		err := response.Success(c, data, "Success message")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := parseResponse(t, rec)
		assert.True(t, resp.Success)
		assert.Equal(t, "Success message", resp.Message)
		assert.NotNil(t, resp.Data)
	})

	t.Run("returns 200 with nil data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.Success(c, nil, "")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		resp := parseResponse(t, rec)
		assert.True(t, resp.Success)
	})

	t.Run("returns 200 with empty message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.Success(c, "test", "")

		assert.NoError(t, err)
		resp := parseResponse(t, rec)
		assert.Equal(t, "", resp.Message)
	})
}

// =============================================================================
// Created Response Tests
// =============================================================================

func TestCreated(t *testing.T) {
	e := echo.New()

	t.Run("returns 201 with data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		data := map[string]string{"id": "123"}
		err := response.Created(c, data, "Created successfully")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		resp := parseResponse(t, rec)
		assert.True(t, resp.Success)
		assert.Equal(t, "Created successfully", resp.Message)
	})

	t.Run("returns 201 with nil data", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.Created(c, nil, "Resource created")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

// =============================================================================
// NoContent Response Tests
// =============================================================================

func TestNoContent(t *testing.T) {
	e := echo.New()

	t.Run("returns 204 with no body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.NoContent(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())
	})
}

// =============================================================================
// Paginated Response Tests
// =============================================================================

func TestPaginated(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name          string
		totalCount    int64
		page          int
		pageSize      int
		expectedPages int
	}{
		{
			name:          "exact division",
			totalCount:    100,
			page:          1,
			pageSize:      10,
			expectedPages: 10,
		},
		{
			name:          "with remainder",
			totalCount:    105,
			page:          1,
			pageSize:      10,
			expectedPages: 11,
		},
		{
			name:          "single page",
			totalCount:    5,
			page:          1,
			pageSize:      10,
			expectedPages: 1,
		},
		{
			name:          "empty result",
			totalCount:    0,
			page:          1,
			pageSize:      10,
			expectedPages: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			data := []string{"item1", "item2"}
			err := response.Paginated(c, data, tt.totalCount, tt.page, tt.pageSize)

			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)

			resp := parseResponse(t, rec)
			assert.True(t, resp.Success)
			require.NotNil(t, resp.Meta)
			assert.Equal(t, tt.page, resp.Meta.Page)
			assert.Equal(t, tt.pageSize, resp.Meta.PageSize)
			assert.Equal(t, tt.totalCount, resp.Meta.TotalCount)
			assert.Equal(t, tt.expectedPages, resp.Meta.TotalPages)
		})
	}
}

// =============================================================================
// Error Response Tests
// =============================================================================

func TestError(t *testing.T) {
	e := echo.New()

	t.Run("returns error response", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.Error(c, http.StatusBadRequest, "CUSTOM_ERROR", "Custom error message")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		require.NotNil(t, resp.Error)
		assert.Equal(t, "CUSTOM_ERROR", resp.Error.Code)
		assert.Equal(t, "Custom error message", resp.Error.Message)
	})
}

func TestErrorWithDetails(t *testing.T) {
	e := echo.New()

	t.Run("returns error with details", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		details := map[string]string{
			"field1": "error 1",
			"field2": "error 2",
		}
		err := response.ErrorWithDetails(c, http.StatusBadRequest, "VALIDATION", "Validation failed", details)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		require.NotNil(t, resp.Error)
		assert.Equal(t, "VALIDATION", resp.Error.Code)
		assert.Equal(t, "Validation failed", resp.Error.Message)
		assert.Equal(t, "error 1", resp.Error.Details["field1"])
		assert.Equal(t, "error 2", resp.Error.Details["field2"])
	})
}

// =============================================================================
// Convenience Error Functions Tests
// =============================================================================

func TestBadRequest(t *testing.T) {
	e := echo.New()

	t.Run("returns 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.BadRequest(c, "Invalid input")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		assert.Equal(t, "BAD_REQUEST", resp.Error.Code)
	})
}

func TestUnauthorized(t *testing.T) {
	e := echo.New()

	t.Run("returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.Unauthorized(c, "Invalid token")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		assert.Equal(t, "UNAUTHORIZED", resp.Error.Code)
	})
}

func TestForbidden(t *testing.T) {
	e := echo.New()

	t.Run("returns 403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.Forbidden(c, "Access denied")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		assert.Equal(t, "FORBIDDEN", resp.Error.Code)
	})
}

func TestNotFound(t *testing.T) {
	e := echo.New()

	t.Run("returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.NotFound(c, "Resource not found")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		assert.Equal(t, "NOT_FOUND", resp.Error.Code)
	})
}

func TestConflict(t *testing.T) {
	e := echo.New()

	t.Run("returns 409", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.Conflict(c, "Resource already exists")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		assert.Equal(t, "CONFLICT", resp.Error.Code)
	})
}

func TestInternalError(t *testing.T) {
	e := echo.New()

	t.Run("returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := response.InternalError(c, "Something went wrong")

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		assert.Equal(t, "INTERNAL_ERROR", resp.Error.Code)
	})
}

func TestValidationError(t *testing.T) {
	e := echo.New()

	t.Run("returns 400 with details", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		details := map[string]string{
			"name":  "Name is required",
			"email": "Invalid email format",
		}
		err := response.ValidationError(c, details)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp := parseResponse(t, rec)
		assert.False(t, resp.Success)
		assert.Equal(t, "VALIDATION_ERROR", resp.Error.Code)
		assert.Equal(t, "Validation failed", resp.Error.Message)
		assert.Len(t, resp.Error.Details, 2)
	})
}

// =============================================================================
// Content-Type Tests
// =============================================================================

func TestResponse_ContentType(t *testing.T) {
	e := echo.New()

	t.Run("returns JSON content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		_ = response.Success(c, nil, "")

		assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	})
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkSuccess(b *testing.B) {
	e := echo.New()
	data := map[string]string{"key": "value"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = response.Success(c, data, "Success")
	}
}

func BenchmarkPaginated(b *testing.B) {
	e := echo.New()
	data := []string{"item1", "item2", "item3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = response.Paginated(c, data, 100, 1, 10)
	}
}

func BenchmarkError(b *testing.B) {
	e := echo.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = response.BadRequest(c, "Error message")
	}
}
