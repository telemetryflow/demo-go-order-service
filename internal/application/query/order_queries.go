// Package query contains CQRS queries for Order.
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

// GetAllordersQuery represents the get all orders query with pagination
type GetAllordersQuery struct {
	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}

// Validate validates the query
func (q *GetAllordersQuery) Validate() error {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	return nil
}

// SearchordersQuery represents the search orders query
type SearchordersQuery struct {
	Query  string `json:"query" query:"query"`
	Offset int    `json:"offset" query:"offset"`
	Limit  int    `json:"limit" query:"limit"`
}

// Validate validates the query
func (q *SearchordersQuery) Validate() error {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	return nil
}
