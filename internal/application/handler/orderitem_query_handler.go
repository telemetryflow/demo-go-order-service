// Package handler handles CQRS queries for OrderItem.
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

package handler

import (
	"context"

	"github.com/telemetryflow/order-service/internal/application/dto"
	"github.com/telemetryflow/order-service/internal/application/query"
	"github.com/telemetryflow/order-service/internal/domain/repository"
)

// OrderitemQueryHandler handles queries for Orderitem entity
type OrderitemQueryHandler struct {
	repo repository.OrderitemRepository
}

// NewOrderitemQueryHandler creates a new Orderitem query handler
func NewOrderitemQueryHandler(repo repository.OrderitemRepository) *OrderitemQueryHandler {
	return &OrderitemQueryHandler{
		repo: repo,
	}
}

// HandleOrderitemGetByID handles get orderitem by ID query
func (h *OrderitemQueryHandler) HandleOrderitemGetByID(ctx context.Context, qry *query.GetOrderitemByIDQuery) (*dto.OrderitemResponse, error) {
	entity, err := h.repo.FindByID(ctx, qry.ID)
	if err != nil {
		return nil, err
	}
	return dto.OrderitemToResponse(entity), nil
}

// HandleOrderitemGetAll handles get all orderitems query
func (h *OrderitemQueryHandler) HandleOrderitemGetAll(ctx context.Context, qry *query.GetAllOrderItemsQuery) (*dto.OrderitemListResponse, error) {
	entities, total, err := h.repo.FindAll(ctx, qry.Offset, qry.Limit)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.OrderitemResponse, len(entities))
	for i := range entities {
		responses[i] = dto.OrderitemToResponse(&entities[i])
	}

	return &dto.OrderitemListResponse{
		Data:   responses,
		Total:  int(total),
		Offset: qry.Offset,
		Limit:  qry.Limit,
	}, nil
}
