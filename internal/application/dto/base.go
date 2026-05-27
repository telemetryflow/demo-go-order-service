// Package dto contains base Data Transfer Object types.
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

package dto

import (
	"time"

	"github.com/google/uuid"
)

// BaseResponse contains common response fields
type BaseResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to the error response
func (e ErrorResponse) WithDetails(details map[string]string) ErrorResponse {
	e.Details = details
	return e
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}, message string) SuccessResponse {
	return SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// PaginatedResponse represents a paginated response
type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](items []T, totalCount int64, page, pageSize int) PaginatedResponse[T] {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}
	return PaginatedResponse[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
