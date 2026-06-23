// Package domain provides domain layer types and aggregates.
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

package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common domain errors
var (
	ErrEntityNotFound = errors.New("entity not found")
	ErrInvalidEntity  = errors.New("invalid entity")
)

// Entity represents a generic domain entity for testing purposes
type Entity struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewEntity creates a new Entity with the given ID and name
func NewEntity(id uuid.UUID, name string) *Entity {
	now := time.Now()
	return &Entity{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewEntityWithValidation creates a new Entity with validation
func NewEntityWithValidation(id uuid.UUID, name string) (*Entity, error) {
	if err := ValidateEntityName(name); err != nil {
		return nil, err
	}
	return NewEntity(id, name), nil
}

// ValidateEntityName validates an entity name
func ValidateEntityName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if len(name) > 255 {
		return errors.New("name exceeds maximum length")
	}
	return nil
}
