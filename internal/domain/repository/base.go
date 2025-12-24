// Package repository defines repository interfaces.
package repository

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the base repository interface
type Repository[T any] interface {
	// Create creates a new entity
	Create(ctx context.Context, entity *T) error

	// FindByID finds an entity by ID
	FindByID(ctx context.Context, id uuid.UUID) (*T, error)

	// FindAll finds all entities with pagination
	FindAll(ctx context.Context, offset, limit int) ([]T, int64, error)

	// Update updates an existing entity
	Update(ctx context.Context, entity *T) error

	// Delete deletes an entity by ID (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// HardDelete permanently deletes an entity
	HardDelete(ctx context.Context, id uuid.UUID) error
}

// Pagination holds pagination parameters
type Pagination struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"page_size" query:"page_size"`
}

// NewPagination creates pagination with defaults
func NewPagination(page, pageSize int) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset calculates the offset for database queries
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the limit for database queries
func (p Pagination) Limit() int {
	return p.PageSize
}

// PaginatedResult holds paginated query results
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginatedResult creates a new paginated result
func NewPaginatedResult[T any](items []T, totalCount int64, page, pageSize int) PaginatedResult[T] {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}
	return PaginatedResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
