// Package handler handles CQRS commands for OrderItem.
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

// OrderitemCommandHandler handles commands for Orderitem entity
type OrderitemCommandHandler struct {
	repo repository.OrderitemRepository
}

// NewOrderitemCommandHandler creates a new Orderitem command handler
func NewOrderitemCommandHandler(repo repository.OrderitemRepository) *OrderitemCommandHandler {
	return &OrderitemCommandHandler{
		repo: repo,
	}
}

// HandleOrderitemCreate handles create orderitem command
func (h *OrderitemCommandHandler) HandleOrderitemCreate(ctx context.Context, cmd *command.CreateOrderitemCommand) error {
	entity := cmd.ToEntity()
	return h.repo.Create(ctx, entity)
}

// HandleOrderitemUpdate handles update orderitem command
func (h *OrderitemCommandHandler) HandleOrderitemUpdate(ctx context.Context, cmd *command.UpdateOrderitemCommand) error {
	entity := cmd.ToEntity()
	return h.repo.Update(ctx, entity)
}

// HandleOrderitemDelete handles delete orderitem command
func (h *OrderitemCommandHandler) HandleOrderitemDelete(ctx context.Context, cmd *command.DeleteOrderitemCommand) error {
	return h.repo.Delete(ctx, cmd.ID)
}
