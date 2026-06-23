// Package handler handles CQRS commands for Order.
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

	"github.com/telemetryflow/order-service/internal/application/command"
	"github.com/telemetryflow/order-service/internal/domain/repository"
)

// OrderCommandHandler handles commands for Order entity
type OrderCommandHandler struct {
	repo repository.OrderRepository
}

// NewOrderCommandHandler creates a new Order command handler
func NewOrderCommandHandler(repo repository.OrderRepository) *OrderCommandHandler {
	return &OrderCommandHandler{
		repo: repo,
	}
}

// HandleOrderCreate handles create order command
func (h *OrderCommandHandler) HandleOrderCreate(ctx context.Context, cmd *command.CreateOrderCommand) error {
	entity := cmd.ToEntity()
	return h.repo.Create(ctx, entity)
}

// HandleOrderUpdate handles update order command
func (h *OrderCommandHandler) HandleOrderUpdate(ctx context.Context, cmd *command.UpdateOrderCommand) error {
	entity := cmd.ToEntity()
	return h.repo.Update(ctx, entity)
}

// HandleOrderDelete handles delete order command
func (h *OrderCommandHandler) HandleOrderDelete(ctx context.Context, cmd *command.DeleteOrderCommand) error {
	return h.repo.Delete(ctx, cmd.ID)
}
