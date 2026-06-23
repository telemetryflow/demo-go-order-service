// Package entity contains base domain entity types.
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

package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common fields for all entities
type Base struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// NewBase creates a new Base with generated ID and timestamps
func NewBase() Base {
	now := time.Now()
	return Base{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsDeleted returns true if the entity has been soft-deleted
func (b *Base) IsDeleted() bool {
	return b.DeletedAt.Valid
}

// MarkUpdated updates the UpdatedAt timestamp
func (b *Base) MarkUpdated() {
	b.UpdatedAt = time.Now()
}

// MarkDeleted sets the DeletedAt timestamp for soft delete
func (b *Base) MarkDeleted() {
	now := time.Now()
	b.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
}

// Restore clears the DeletedAt timestamp
func (b *Base) Restore() {
	b.DeletedAt = gorm.DeletedAt{}
}

// BeforeCreate is a GORM hook that runs before creating a record
func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
