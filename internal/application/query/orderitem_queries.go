// Package query contains CQRS queries for Orderitem.
package query

import (
	"github.com/google/uuid"
)

// GetOrderitemByIDQuery represents the get orderitem by ID query
type GetOrderitemByIDQuery struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// Validate validates the query
func (q *GetOrderitemByIDQuery) Validate() error {
	if q.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}

// ListOrderitemsQuery represents the list orderitems query
type ListOrderitemsQuery struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	SortDir  string `json:"sort_dir" query:"sort_dir"`
	Search   string `json:"search" query:"search"`
}

// Validate validates the query
func (q *ListOrderitemsQuery) Validate() error {
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
func (q *ListOrderitemsQuery) Offset() int {
	return (q.Page - 1) * q.PageSize
}

// GetAllorderitemsQuery represents the get all orderitems query with pagination
type GetAllorderitemsQuery struct {
	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}

// Validate validates the query
func (q *GetAllorderitemsQuery) Validate() error {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	return nil
}

// SearchorderitemsQuery represents the search orderitems query
type SearchorderitemsQuery struct {
	Query  string `json:"query" query:"query"`
	Offset int    `json:"offset" query:"offset"`
	Limit  int    `json:"limit" query:"limit"`
}

// Validate validates the query
func (q *SearchorderitemsQuery) Validate() error {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	return nil
}
