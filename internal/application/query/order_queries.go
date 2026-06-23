// Package query contains CQRS queries for Order.
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
	"github.com/google/uuid"
)

// GetOrderByIDQuery represents the get order by ID query
type GetOrderByIDQuery struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// Validate validates the query
func (q *GetOrderByIDQuery) Validate() error {
	if q.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}

// ListOrdersQuery represents the list orders query
type ListOrdersQuery struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	SortDir  string `json:"sort_dir" query:"sort_dir"`
	Search   string `json:"search" query:"search"`
}

// Validate validates the query
func (q *ListOrdersQuery) Validate() error {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 10
	}
	if q.SortDir != "asc" && q.SortDir != "desc" {
		q.SortDir = "desc"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	return nil
}

// Offset returns the offset for pagination
func (q *ListOrdersQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

// GetAllOrdersQuery represents the get all orders query with pagination
type GetAllOrdersQuery struct {
	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}

// Validate validates the query
func (q *GetAllOrdersQuery) Validate() error {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	return nil
}

// SearchOrdersQuery represents the search orders query
type SearchOrdersQuery struct {
	Query  string `json:"query" query:"query"`
	Offset int    `json:"offset" query:"offset"`
	Limit  int    `json:"limit" query:"limit"`
}

// Validate validates the query
func (q *SearchOrdersQuery) Validate() error {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	return nil
}
