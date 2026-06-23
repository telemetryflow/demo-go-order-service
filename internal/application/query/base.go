// Package query contains base CQRS query types.
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

package query

import (
	"context"

	"github.com/google/uuid"
)

// Query represents a read operation
type Query interface {
	Validate() error
}

// QueryHandler handles query execution
type QueryHandler[T any] interface {
	Handle(ctx context.Context, q Query) (T, error)
}

// GetByIDQuery is a base for get-by-id operations
type GetByIDQuery struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// Validate validates the get-by-id query
func (q *GetByIDQuery) Validate() error {
	if q.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}

// ListQuery is a base for list operations
type ListQuery struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	SortDir  string `json:"sort_dir" query:"sort_dir"`
	Search   string `json:"search" query:"search"`
}

// Validate validates the list query
func (q *ListQuery) Validate() error {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 10
	}
	if q.SortDir != "asc" && q.SortDir != "desc" {
		q.SortDir = "desc"
	}
	return nil
}

// Offset returns the offset for pagination
func (q *ListQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

// ListResult holds paginated list results
type ListResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewListResult creates a new list result
func NewListResult[T any](items []T, totalCount int64, page, pageSize int) ListResult[T] {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}
	return ListResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// Common query errors
var (
	ErrInvalidID = &QueryError{Code: "INVALID_ID", Message: "Invalid ID provided"}
	ErrNotFound  = &QueryError{Code: "NOT_FOUND", Message: "Resource not found"}
	ErrForbidden = &QueryError{Code: "FORBIDDEN", Message: "Access forbidden"}
)

// QueryError represents a query execution error
type QueryError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *QueryError) Error() string {
	return e.Message
}
