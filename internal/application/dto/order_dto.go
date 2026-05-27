// Package dto contains DTOs for Order endpoints.
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
	"github.com/telemetryflow/order-service/internal/domain/entity"
)

// OrderResponse represents the order API response
type OrderResponse struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Total      float64   `json:"total"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FromOrder converts entity to response DTO
func FromOrder(e *entity.Order) OrderResponse {
	return OrderResponse{
		ID:         e.ID,
		CustomerID: e.CustomerID,
		Total:      e.Total,
		Status:     e.Status,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

// FromOrders converts entities to response DTOs
func FromOrders(entities []entity.Order) []OrderResponse {
	responses := make([]OrderResponse, len(entities))
	for i, e := range entities {
		responses[i] = FromOrder(&e)
	}
	return responses
}

// CreateOrderRequest represents the create order request
type CreateOrderRequest struct {
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Total      float64   `json:"total" validate:"required"`
	Status     string    `json:"status" validate:"required"`
}

// UpdateOrderRequest represents the update order request
type UpdateOrderRequest struct {
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Total      float64   `json:"total" validate:"required"`
	Status     string    `json:"status" validate:"required"`
}

// OrderToResponse converts entity pointer to response DTO pointer
func OrderToResponse(e *entity.Order) *OrderResponse {
	if e == nil {
		return nil
	}
	resp := FromOrder(e)
	return &resp
}

// OrderListResponse represents the list order API response
type OrderListResponse struct {
	Data   []*OrderResponse `json:"data"`
	Total  int              `json:"total"`
	Offset int              `json:"offset"`
	Limit  int              `json:"limit"`
}
