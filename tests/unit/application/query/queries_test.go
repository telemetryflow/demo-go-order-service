// Package query_test provides the query_test package for TelemetryFlow Order Service.
//
// TelemetryFlow Order Service - AI-Powered Observability & Incident Response Management (IRM) Platform
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

package query_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/telemetryflow/order-service/internal/application/query"
)

// =============================================================================
// GetOrderByIDQuery Tests
//
// Tests for querying a single order by its UUID identifier.
// =============================================================================

// TestGetOrderByIDQuery_Validate verifies ID validation rules.
func TestGetOrderByIDQuery_Validate(t *testing.T) {
	tests := []struct {
		name        string
		query       *query.GetOrderByIDQuery
		expectError bool
	}{
		{
			name:        "valid UUID",
			query:       &query.GetOrderByIDQuery{ID: uuid.New()},
			expectError: false,
		},
		{
			name:        "nil UUID",
			query:       &query.GetOrderByIDQuery{ID: uuid.Nil},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, query.ErrInvalidID, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// =============================================================================
// ListOrdersQuery Tests
//
// Tests for paginated order listing with sorting and page size validation.
// =============================================================================

// TestListOrdersQuery_Validate verifies pagination and sorting validation.
func TestListOrdersQuery_Validate(t *testing.T) {
	tests := []struct {
		name             string
		query            *query.ListOrdersQuery
		expectedPage     int
		expectedPageSize int
		expectedSortDir  string
		expectedSortBy   string
	}{
		{
			name: "valid query",
			query: &query.ListOrdersQuery{
				Page:     1,
				PageSize: 20,
				SortDir:  "asc",
				SortBy:   "total",
			},
			expectedPage:     1,
			expectedPageSize: 20,
			expectedSortDir:  "asc",
			expectedSortBy:   "total",
		},
		{
			name:             "defaults for invalid page",
			query:            &query.ListOrdersQuery{Page: 0},
			expectedPage:     1,
			expectedPageSize: 10,
			expectedSortDir:  "desc",
			expectedSortBy:   "created_at",
		},
		{
			name:             "defaults for negative page",
			query:            &query.ListOrdersQuery{Page: -5},
			expectedPage:     1,
			expectedPageSize: 10,
			expectedSortDir:  "desc",
			expectedSortBy:   "created_at",
		},
		{
			name:             "defaults for invalid page size (0)",
			query:            &query.ListOrdersQuery{PageSize: 0},
			expectedPage:     1,
			expectedPageSize: 10,
			expectedSortDir:  "desc",
			expectedSortBy:   "created_at",
		},
		{
			name:             "defaults for page size over limit",
			query:            &query.ListOrdersQuery{PageSize: 200},
			expectedPage:     1,
			expectedPageSize: 10,
			expectedSortDir:  "desc",
			expectedSortBy:   "created_at",
		},
		{
			name:             "defaults for invalid sort direction",
			query:            &query.ListOrdersQuery{SortDir: "invalid"},
			expectedPage:     1,
			expectedPageSize: 10,
			expectedSortDir:  "desc",
			expectedSortBy:   "created_at",
		},
		{
			name:             "defaults for empty sort by",
			query:            &query.ListOrdersQuery{SortBy: ""},
			expectedPage:     1,
			expectedPageSize: 10,
			expectedSortDir:  "desc",
			expectedSortBy:   "created_at",
		},
		{
			name: "asc sort direction",
			query: &query.ListOrdersQuery{
				Page:     1,
				PageSize: 10,
				SortDir:  "asc",
			},
			expectedPage:     1,
			expectedPageSize: 10,
			expectedSortDir:  "asc",
			expectedSortBy:   "created_at",
		},
		{
			name: "desc sort direction",
			query: &query.ListOrdersQuery{
				Page:     1,
				PageSize: 10,
				SortDir:  "desc",
			},
			expectedPage:     1,
			expectedPageSize: 10,
			expectedSortDir:  "desc",
			expectedSortBy:   "created_at",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedPage, tt.query.Page)
			assert.Equal(t, tt.expectedPageSize, tt.query.PageSize)
			assert.Equal(t, tt.expectedSortDir, tt.query.SortDir)
			assert.Equal(t, tt.expectedSortBy, tt.query.SortBy)
		})
	}
}

func TestListOrdersQuery_Offset(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		expectedOffset int
	}{
		{
			name:           "first page",
			page:           1,
			pageSize:       10,
			expectedOffset: 0,
		},
		{
			name:           "second page",
			page:           2,
			pageSize:       10,
			expectedOffset: 10,
		},
		{
			name:           "third page",
			page:           3,
			pageSize:       10,
			expectedOffset: 20,
		},
		{
			name:           "custom page size",
			page:           5,
			pageSize:       25,
			expectedOffset: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &query.ListOrdersQuery{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}
			assert.Equal(t, tt.expectedOffset, q.Offset())
		})
	}
}

// =============================================================================
// GetAllOrdersQuery Tests
//
// Tests for offset/limit based order retrieval with validation.
// =============================================================================

// TestGetAllOrdersQuery_Validate verifies offset and limit validation.
func TestGetAllOrdersQuery_Validate(t *testing.T) {
	tests := []struct {
		name           string
		query          *query.GetAllOrdersQuery
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "valid query",
			query:          &query.GetAllOrdersQuery{Offset: 10, Limit: 20},
			expectedOffset: 10,
			expectedLimit:  20,
		},
		{
			name:           "defaults for negative offset",
			query:          &query.GetAllOrdersQuery{Offset: -5, Limit: 10},
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "defaults for invalid limit (0)",
			query:          &query.GetAllOrdersQuery{Offset: 0, Limit: 0},
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "defaults for negative limit",
			query:          &query.GetAllOrdersQuery{Offset: 0, Limit: -10},
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "defaults for limit over 100",
			query:          &query.GetAllOrdersQuery{Offset: 0, Limit: 200},
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "maximum allowed limit",
			query:          &query.GetAllOrdersQuery{Offset: 0, Limit: 100},
			expectedOffset: 0,
			expectedLimit:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOffset, tt.query.Offset)
			assert.Equal(t, tt.expectedLimit, tt.query.Limit)
		})
	}
}

// =============================================================================
// SearchOrdersQuery Tests
//
// Tests for full-text search with pagination and validation.
// =============================================================================

// TestSearchOrdersQuery_Validate verifies search query validation.
func TestSearchOrdersQuery_Validate(t *testing.T) {
	tests := []struct {
		name           string
		query          *query.SearchOrdersQuery
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "valid query",
			query:          &query.SearchOrdersQuery{Query: "search term", Offset: 10, Limit: 20},
			expectedOffset: 10,
			expectedLimit:  20,
		},
		{
			name:           "defaults for negative offset",
			query:          &query.SearchOrdersQuery{Query: "test", Offset: -5, Limit: 10},
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "defaults for invalid limit",
			query:          &query.SearchOrdersQuery{Query: "test", Offset: 0, Limit: 0},
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "defaults for limit over 100",
			query:          &query.SearchOrdersQuery{Query: "test", Offset: 0, Limit: 150},
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "empty search query",
			query:          &query.SearchOrdersQuery{Query: "", Offset: 0, Limit: 10},
			expectedOffset: 0,
			expectedLimit:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.query.Validate()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOffset, tt.query.Offset)
			assert.Equal(t, tt.expectedLimit, tt.query.Limit)
		})
	}
}

// =============================================================================
// Edge Cases
//
// Tests for unusual query inputs and boundary conditions.
// =============================================================================

// TestQuery_EdgeCases verifies query behavior under edge conditions.
func TestQuery_EdgeCases(t *testing.T) {
	t.Run("ListOrdersQuery with max page", func(t *testing.T) {
		q := &query.ListOrdersQuery{
			Page:     10000,
			PageSize: 100,
		}
		_ = q.Validate()
		assert.Equal(t, 10000, q.Page)
		// PageSize 100 is valid (limit is <= 100), so it should be reset to default
		assert.Equal(t, 100, q.PageSize)
	})

	t.Run("GetOrderByIDQuery with specific UUID", func(t *testing.T) {
		specificID := uuid.MustParse("12345678-1234-1234-1234-123456789012")
		q := &query.GetOrderByIDQuery{ID: specificID}
		err := q.Validate()
		assert.NoError(t, err)
		assert.Equal(t, specificID, q.ID)
	})

	t.Run("SearchOrdersQuery with special characters", func(t *testing.T) {
		q := &query.SearchOrdersQuery{
			Query:  "test@#$%^&*()",
			Offset: 0,
			Limit:  10,
		}
		err := q.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "test@#$%^&*()", q.Query)
	})
}

// =============================================================================
// Benchmark Tests
//
// Performance benchmarks for query operations.
// Run with: go test -bench=. -benchmem
// =============================================================================

// BenchmarkListOrdersQuery_Validate measures validation performance.
func BenchmarkListOrdersQuery_Validate(b *testing.B) {
	q := &query.ListOrdersQuery{
		Page:     1,
		PageSize: 20,
		SortDir:  "asc",
		SortBy:   "created_at",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Validate()
	}
}

func BenchmarkListOrdersQuery_Offset(b *testing.B) {
	q := &query.ListOrdersQuery{
		Page:     100,
		PageSize: 50,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Offset()
	}
}

func BenchmarkGetOrderByIDQuery_Validate(b *testing.B) {
	q := &query.GetOrderByIDQuery{ID: uuid.New()}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Validate()
	}
}
