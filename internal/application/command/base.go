// Package command contains base CQRS command types.
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

package command

import (
	"context"

	"github.com/google/uuid"
)

// Command represents a write operation
type Command interface {
	Validate() error
}

// CommandHandler handles command execution
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// CreateCommand is a base for create operations
type CreateCommand struct {
	// Embed entity-specific fields
}

// UpdateCommand is a base for update operations
type UpdateCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// DeleteCommand is a base for delete operations
type DeleteCommand struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// Validate validates the delete command
func (c *DeleteCommand) Validate() error {
	if c.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}

// CommandResult represents the result of a command execution
type CommandResult struct {
	ID      uuid.UUID `json:"id,omitempty"`
	Success bool      `json:"success"`
	Message string    `json:"message,omitempty"`
}

// NewSuccessResult creates a success result
func NewSuccessResult(id uuid.UUID, message string) CommandResult {
	return CommandResult{
		ID:      id,
		Success: true,
		Message: message,
	}
}

// NewErrorResult creates an error result
func NewErrorResult(message string) CommandResult {
	return CommandResult{
		Success: false,
		Message: message,
	}
}

// Common command errors
var (
	ErrInvalidID     = &CommandError{Code: "INVALID_ID", Message: "Invalid ID provided"}
	ErrValidation    = &CommandError{Code: "VALIDATION_ERROR", Message: "Validation failed"}
	ErrNotFound      = &CommandError{Code: "NOT_FOUND", Message: "Resource not found"}
	ErrAlreadyExists = &CommandError{Code: "ALREADY_EXISTS", Message: "Resource already exists"}
	ErrUnauthorized  = &CommandError{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
)

// CommandError represents a command execution error
type CommandError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *CommandError) Error() string {
	return e.Message
}
