// Package handler handles CQRS queries for Order.
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

package handler

import (
	"context"

	"github.com/telemetryflow/order-service/internal/application/dto"
	"github.com/telemetryflow/order-service/internal/application/query"
	"github.com/telemetryflow/order-service/internal/domain/repository"
)

// OrderQueryHandler handles queries for Order entity
type OrderQueryHandler struct {
	repo repository.OrderRepository
}

// NewOrderQueryHandler creates a new Order query handler
func NewOrderQueryHandler(repo repository.OrderRepository) *OrderQueryHandler {
	return &OrderQueryHandler{
		repo: repo,
	}
}

// HandleOrderGetByID handles get order by ID query
func (h *OrderQueryHandler) HandleOrderGetByID(ctx context.Context, qry *query.GetOrderByIDQuery) (*dto.OrderResponse, error) {
	entity, err := h.repo.FindByID(ctx, qry.ID)
	if err != nil {
		return nil, err
	}
	return dto.OrderToResponse(entity), nil
}

// HandleOrderGetAll handles get all orders query
func (h *OrderQueryHandler) HandleOrderGetAll(ctx context.Context, qry *query.GetAllOrdersQuery) (*dto.OrderListResponse, error) {
	entities, total, err := h.repo.FindAll(ctx, qry.Offset, qry.Limit)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.OrderResponse, len(entities))
	for i := range entities {
		responses[i] = dto.OrderToResponse(&entities[i])
	}

	return &dto.OrderListResponse{
		Data:   responses,
		Total:  int(total),
		Offset: qry.Offset,
		Limit:  qry.Limit,
	}, nil
}
