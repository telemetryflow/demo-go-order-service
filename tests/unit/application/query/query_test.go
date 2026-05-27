// Package query_test provides unit tests for CQRS queries.
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

package query_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/order-service/internal/application/query"
)

// =============================================================================
// Order Query Tests
// =============================================================================

func TestGetOrderByIDQuery(t *testing.T) {
	t.Run("should validate with valid ID", func(t *testing.T) {
		q := &query.GetOrderByIDQuery{
			ID: uuid.New(),
		}
		err := q.Validate()
		require.NoError(t, err)
	})

	t.Run("should fail validation with nil ID", func(t *testing.T) {
		q := &query.GetOrderByIDQuery{
			ID: uuid.Nil,
		}
		err := q.Validate()
		require.Error(t, err)
	})
}

func TestGetAllOrdersQuery(t *testing.T) {
	t.Run("should set default values for invalid offset", func(t *testing.T) {
		q := &query.GetAllOrdersQuery{
			Offset: -1,
			Limit:  10,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 0, q.Offset)
	})

	t.Run("should set default limit for invalid limit", func(t *testing.T) {
		q := &query.GetAllOrdersQuery{
			Offset: 0,
			Limit:  0,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 10, q.Limit)
	})

	t.Run("should cap limit at 100", func(t *testing.T) {
		q := &query.GetAllOrdersQuery{
			Offset: 0,
			Limit:  200,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 10, q.Limit)
	})

	t.Run("should accept valid pagination", func(t *testing.T) {
		q := &query.GetAllOrdersQuery{
			Offset: 20,
			Limit:  50,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 20, q.Offset)
		assert.Equal(t, 50, q.Limit)
	})
}

func TestListOrdersQuery(t *testing.T) {
	t.Run("should set default page for invalid page", func(t *testing.T) {
		q := &query.ListOrdersQuery{
			Page:     0,
			PageSize: 10,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 1, q.Page)
	})

	t.Run("should set default page size for invalid size", func(t *testing.T) {
		q := &query.ListOrdersQuery{
			Page:     1,
			PageSize: 0,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 10, q.PageSize)
	})

	t.Run("should set default sort direction", func(t *testing.T) {
		q := &query.ListOrdersQuery{
			Page:     1,
			PageSize: 10,
			SortDir:  "invalid",
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, "desc", q.SortDir)
	})

	t.Run("should set default sort by", func(t *testing.T) {
		q := &query.ListOrdersQuery{
			Page:     1,
			PageSize: 10,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, "created_at", q.SortBy)
	})

	t.Run("should calculate offset correctly", func(t *testing.T) {
		q := &query.ListOrdersQuery{
			Page:     3,
			PageSize: 10,
		}
		assert.Equal(t, 20, q.Offset())
	})
}

func TestSearchOrdersQuery(t *testing.T) {
	t.Run("should set default offset for invalid offset", func(t *testing.T) {
		q := &query.SearchOrdersQuery{
			Query:  "test",
			Offset: -1,
			Limit:  10,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 0, q.Offset)
	})

	t.Run("should set default limit for invalid limit", func(t *testing.T) {
		q := &query.SearchOrdersQuery{
			Query:  "test",
			Offset: 0,
			Limit:  0,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 10, q.Limit)
	})
}

// =============================================================================
// Order Item Query Tests
// =============================================================================

func TestGetOrderitemByIDQuery(t *testing.T) {
	t.Run("should validate with valid ID", func(t *testing.T) {
		q := &query.GetOrderitemByIDQuery{
			ID: uuid.New(),
		}
		err := q.Validate()
		require.NoError(t, err)
	})

	t.Run("should fail validation with nil ID", func(t *testing.T) {
		q := &query.GetOrderitemByIDQuery{
			ID: uuid.Nil,
		}
		err := q.Validate()
		require.Error(t, err)
	})
}

func TestGetAllOrderItemsQuery(t *testing.T) {
	t.Run("should set default values for invalid offset", func(t *testing.T) {
		q := &query.GetAllOrderItemsQuery{
			Offset: -1,
			Limit:  10,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 0, q.Offset)
	})

	t.Run("should set default limit for invalid limit", func(t *testing.T) {
		q := &query.GetAllOrderItemsQuery{
			Offset: 0,
			Limit:  0,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 10, q.Limit)
	})

	t.Run("should cap limit at 100", func(t *testing.T) {
		q := &query.GetAllOrderItemsQuery{
			Offset: 0,
			Limit:  200,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 10, q.Limit)
	})

	t.Run("should accept valid pagination", func(t *testing.T) {
		q := &query.GetAllOrderItemsQuery{
			Offset: 20,
			Limit:  50,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 20, q.Offset)
		assert.Equal(t, 50, q.Limit)
	})
}

func TestListOrderitemsQuery(t *testing.T) {
	t.Run("should set default page for invalid page", func(t *testing.T) {
		q := &query.ListOrderitemsQuery{
			Page:     0,
			PageSize: 10,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 1, q.Page)
	})

	t.Run("should set default page size for invalid size", func(t *testing.T) {
		q := &query.ListOrderitemsQuery{
			Page:     1,
			PageSize: 0,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 10, q.PageSize)
	})

	t.Run("should calculate offset correctly", func(t *testing.T) {
		q := &query.ListOrderitemsQuery{
			Page:     3,
			PageSize: 10,
		}
		assert.Equal(t, 20, q.Offset())
	})
}

func TestSearchOrderItemsQuery(t *testing.T) {
	t.Run("should set default offset for invalid offset", func(t *testing.T) {
		q := &query.SearchOrderItemsQuery{
			Query:  "test",
			Offset: -1,
			Limit:  10,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 0, q.Offset)
	})

	t.Run("should set default limit for invalid limit", func(t *testing.T) {
		q := &query.SearchOrderItemsQuery{
			Query:  "test",
			Offset: 0,
			Limit:  0,
		}
		err := q.Validate()
		require.NoError(t, err)
		assert.Equal(t, 10, q.Limit)
	})
}
