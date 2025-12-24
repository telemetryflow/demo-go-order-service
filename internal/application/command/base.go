// Package command contains CQRS commands (write operations).
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
